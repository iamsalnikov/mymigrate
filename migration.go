package mymigrate

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"time"
)

type UpFunc func(db *sql.DB) error
type DownFunc func(db *sql.DB) error

type mig struct {
	name string
	up   UpFunc
	down DownFunc
}

const migrationsTable = "mymigrations"

var (
	// set of project's migrations
	migrations = make(map[string]mig)
	// database connection
	db *sql.DB
	// function to get list of applied migrations
	getApplied = defaultAppliedFunc
	// function to mark migration as aplied
	markApplied = defaultMarkAppliedFunc
	// function to down migrations
	down = defaultDownFunc
)

func defaultAppliedFunc(db *sql.DB) ([]string, error) {
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
}

func defaultMarkAppliedFunc(db *sql.DB, name string) error {
	err := createMigrationsTable(db)
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	query := fmt.Sprintf("INSERT INTO %s (name, time) VALUES (?, ?)", migrationsTable)
	_, err = db.ExecContext(ctx, query, name, time.Now())

	return err
}

func defaultDownFunc(db *sql.DB, names []string) ([]string, error) {
	err := createMigrationsTable(db)
	if err != nil {
		return []string{}, err
	}

	downed := make([]string, 0, len(names))
	for _, name := range names {
		mig, ok := migrations[name]
		if !ok {
			return downed, fmt.Errorf("can't find migration '%s'", name)
		}

		err = mig.down(db)
		if err != nil {
			return downed, err
		}

		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		query := fmt.Sprintf("DELETE FROM %s WHERE name=?", migrationsTable)
		_, err = db.ExecContext(ctx, query, name)
		if err != nil {
			return downed, err
		}

		downed = append(downed, name)
	}

	return downed, nil
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
func Add(name string, up UpFunc, down DownFunc) {
	migrations[name] = mig{
		name: name,
		up:   up,
		down: down,
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
func Apply() ([]string, error) {
	newNames, err := NewNames()
	if err != nil {
		return []string{}, err
	}

	applied := make([]string, 0, len(newNames))
	for _, name := range newNames {
		err = migrations[name].up(db)
		if err != nil {
			return applied, err
		}

		err = markApplied(db, name)
		if err != nil {
			return applied, err
		}

		applied = append(applied, name)
	}

	return applied, nil
}

// datedMigrationName returns dated migration name
func datedMigrationName(name string) string {
	return fmt.Sprintf("%s-%s", time.Now().Format("20060102-150405"), name)
}

// Template func returns a new migration template and the name of the migration
func Template(pkg, name string) (string, string) {
	if len(pkg) == 0 {
		pkg = "migrations"
	}

	template := `package %s

import (
	"database/sql"

	"github.com/iamsalnikov/mymigrate"
)

func init() {
	mymigrate.Add(
		"%s",
		func(db *sql.DB) error {
			// TODO: write UP logic
			return nil
		},
		func(db *sql.DB) error {
			// TODO: write down logic

			return nil
		},
	)
}
`

	name = datedMigrationName(name)
	return fmt.Sprintf(template, pkg, name), name
}

// History func returns chronological history of applied migrations
func History() ([]string, error) {
	return getApplied(db)
}

// Down func reverts particular number of migrations
// Pass 0 as a number to revert all migrations
func Down(number int) ([]string, error) {
	appliedNames, err := getApplied(db)
	if err != nil {
		return []string{}, err
	}

	if len(appliedNames) == 0 {
		return []string{}, nil
	}

	endIndex := number
	if number >= len(appliedNames) || number == 0 {
		endIndex = len(appliedNames)
	}

	namesToDown := appliedNames[:endIndex]
	return down(db, namesToDown)
}
