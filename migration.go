package mymigrate

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"time"
)

type AppliedFunc func(db *sql.DB) ([]string, error)
type MarkAppliedFunc func(db *sql.DB, name string) error
type UpFunc func(db *sql.DB) error
type DownFunc func(db *sql.DB, names []string) error

type mig struct {
	name string
	up   UpFunc
}

const migrationsTable = "mymigrations"

var defaultAppliedFunc = AppliedFunc(func(db *sql.DB) ([]string, error) {
	err := createMigrationsTable(db)
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	query := fmt.Sprintf("SELECT name FROM %s ORDER BY time DESC, name DESC", migrationsTable)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return []string{}, err
	}

	defer rows.Close()

	res := make([]string, 0)
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return []string{}, nil
		}

		res = append(res, name)
	}

	return res, nil
})

var defaultMarkAppliedFunc = MarkAppliedFunc(func(db *sql.DB, name string) error {
	err := createMigrationsTable(db)
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	query := fmt.Sprintf("INSERT INTO %s (name, time) VALUES (?, ?)", migrationsTable)
	_, err = db.ExecContext(ctx, query, name, time.Now())

	return err
})

var defaultDownFunc = DownFunc(func(db *sql.DB, names []string) error {
	return nil
})

var migrations = make(map[string]mig)
var db *sql.DB
var getApplied = defaultAppliedFunc
var markApplied = defaultMarkAppliedFunc
var down = defaultDownFunc

func resetMigrations() {
	migrations = make(map[string]mig)
}

func resetAppliedFunc() {
	getApplied = defaultAppliedFunc
}

func resetMarkAppliedFunc() {
	markApplied = defaultMarkAppliedFunc
}

func resetDownFunc() {
	down = defaultDownFunc
}

func createMigrationsTable(db *sql.DB) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		name VARCHAR(500) NOT NULL unique,
		time timestamp,
		PRIMARY KEY (name)
	) engine=InnoDB`, migrationsTable)

	_, err := db.Exec(query)
	return err
}

// Add adds mig to queue
// Use this function in init()
func Add(name string, up UpFunc) {
	migrations[name] = mig{
		name: name,
		up:   up,
	}
}

// SetDatabase sets a database that we should use for applying migrations
func SetDatabase(database *sql.DB) {
	db = database
}

// NewNames returns names of new migrations
func NewNames() ([]string, error) {
	appliedNames, err := getApplied(db)
	if err != nil {
		return []string{}, err
	}

	applied := map[string]bool{}
	for _, name := range appliedNames {
		applied[name] = true
	}

	result := make([]string, 0)
	for name := range migrations {
		if !applied[name] {
			result = append(result, name)
		}
	}

	sort.Strings(result)

	return result, nil
}

// Apply func applies migrations
func Apply() error {
	newNames, err := NewNames()
	if err != nil {
		return err
	}

	for _, name := range newNames {
		err = migrations[name].up(db)
		if err != nil {
			return err
		}

		err = markApplied(db, name)
		if err != nil {
			return err
		}
	}

	return nil
}

// Template func returns a new migration template
func Template(pkg, name string) string {
	template := `package %s

import (
	"database/sql"

	"github.com/iamsalnikov/mymigrate"
)

func init() {
	mymigrate.Add("%s", migration.UpFunc(func(db *sql.DB) error {
		return nil
	}))
}
`

	name = fmt.Sprintf("%s-%s", time.Now().Format("20060102-150405"), name)
	return fmt.Sprintf(template, pkg, name)
}

// History func returns chronological history of applied migrations
func History() ([]string, error) {
	return getApplied(db)
}

// Down func reverts particular number of migrations
// Pass 0 as a number to revert all migrations
func Down(number int) error {
	appliedNames, err := getApplied(db)
	if err != nil {
		return err
	}

	if len(appliedNames) == 0 {
		return nil
	}

	endIndex := number
	if number >= len(appliedNames) || number == 0 {
		endIndex = len(appliedNames)
	}

	namesToDown := appliedNames[:endIndex]
	return down(db, namesToDown)
}
