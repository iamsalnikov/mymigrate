package mymigrate

import (
	"context"
	"fmt"
	"sort"
	"time"
)

var (
	// set of project's migrations
	migrations = make(map[string]mig)
	// database provider
	dbProvider DbProvider
	// function to get list of applied migrations
	getApplied = defaultAppliedFunc
	// function to mark migration as aplied
	markApplied = defaultMarkAppliedFunc
	// function to down migrations
	down = defaultDownFunc
)

func defaultAppliedFunc(provider DbProvider) ([]string, error) {
	err := provider.CreateMigrationsTable()
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	return provider.GetApplied(ctx)
}

func defaultMarkAppliedFunc(provider DbProvider, name string) error {
	err := provider.CreateMigrationsTable()
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	return provider.MarkApplied(ctx, name, time.Now())
}

func defaultDownFunc(provider DbProvider, names []string) ([]string, error) {
	err := provider.CreateMigrationsTable()
	if err != nil {
		return nil, err
	}

	downed := make([]string, 0, len(names))
	for _, name := range names {
		mig, ok := migrations[name]
		if !ok {
			return downed, fmt.Errorf("can't find migration '%s'", name)
		}

		err = mig.down(provider.GetDb())
		if err != nil {
			return downed, err
		}

		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		err = provider.DeleteApplied(ctx, name)
		if err != nil {
			return downed, err
		}

		downed = append(downed, name)
	}

	return downed, nil
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

// SetDatabaseProvider sets a DbProvider that we should use for applying migrations
func SetDatabaseProvider(provider DbProvider) {
	dbProvider = provider
}

// NewNames returns names of new migrations
func NewNames() ([]string, error) {
	appliedNames, err := getApplied(dbProvider)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	applied := make([]string, 0, len(newNames))
	for _, name := range newNames {
		err = migrations[name].up(dbProvider.GetDb())
		if err != nil {
			return applied, err
		}

		err = markApplied(dbProvider, name)
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
	return getApplied(dbProvider)
}

// Down func reverts particular number of migrations
// Pass 0 as a number to revert all migrations
func Down(number int) ([]string, error) {
	appliedNames, err := getApplied(dbProvider)
	if err != nil {
		return nil, err
	}

	if len(appliedNames) == 0 {
		return []string{}, nil
	}

	endIndex := number
	if number >= len(appliedNames) || number == 0 {
		endIndex = len(appliedNames)
	}

	namesToDown := appliedNames[:endIndex]
	return down(dbProvider, namesToDown)
}
