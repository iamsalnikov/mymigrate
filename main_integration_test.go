// +build integration

package mymigrate

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
)

var dbConn *sql.DB
var connPort string

func TestMain(m *testing.M) {
	var closer func()
	dbConn, closer = initDB()
	code := m.Run()
	closer()
	os.Exit(code)
}

func initDB() (*sql.DB, func()) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalln(err)
	}

	resource, err := pool.Run("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		log.Fatalln(err)
	}

	var conn *sql.DB
	connPort = resource.GetPort("3306/tcp")
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		conn, err = sql.Open("mysql", getConnString("mysql"))
		if err != nil {
			return err
		}
		return conn.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return conn, func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}
}

func getConnString(name string) string {
	return fmt.Sprintf("root:secret@(localhost:%s)/%s", connPort, name)
}

func getDB(name string) *sql.DB {
	_, err := dbConn.Exec(fmt.Sprintf("create database %s character set utf8 collate utf8_general_ci", name))
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := sql.Open("mysql", getConnString(name))
	if err != nil {
		log.Fatalln(err)
	}

	return conn
}

func tableExists(db *sql.DB, name string) bool {
	row := db.QueryRow("SHOW TABLES LIKE '%" + name + "%'")
	var s string
	row.Scan(&s)
	return s == name
}
