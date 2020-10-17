// +build integration

package mymigrate

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDown_CreatesMigrationTable(t *testing.T) {
	resetMigrations()
	db = getDB("test_down_creates_migration_table")
	_ = Down(1)
	if !tableExists(db, "mymigrations") {
		t.Errorf("Can't find table 'mymigrations' after calling Down() on clean DB")
	}
	db = nil
}

func TestDown_ChangesHistory(t *testing.T) {
	stub := func(db *sql.DB) error { return nil }

	resetMigrations()
	db = getDB("test_down_changes_history")
	Add("mig_1", stub, stub)
	Add("mig_2", stub, stub)
	Add("mig_3", stub, stub)
	Add("mig_4", stub, stub)

	err := Apply()
	assert.Nil(t, err, "Unexpected error during apply")

	err = Down(2)
	assert.Nil(t, err, "Unexpected error during down")

	history, err := History()
	assert.Nil(t, err, "Unexpected error during history")

	assert.EqualValues(t, []string{"mig_2", "mig_1"}, history)
}

func TestDown_RevertsMigrations(t *testing.T) {
	type testCase struct {
		migrations        map[string]string
		downNumber        int
		expHistory        []string
		expExistingTables []string
		expDeletedTables  []string
	}

	testCases := map[string]testCase{
		"empty database": {
			migrations:        map[string]string{},
			downNumber:        1,
			expHistory:        []string{},
			expExistingTables: []string{},
			expDeletedTables:  []string{},
		},
		"three applied migrations and one reverted": {
			migrations: map[string]string{
				"mig_001": "table_001",
				"mig_002": "table_002",
				"mig_003": "table_003",
			},
			downNumber:        1,
			expHistory:        []string{"mig_002", "mig_001"},
			expExistingTables: []string{"table_001", "table_002"},
			expDeletedTables:  []string{"table_003"},
		},
		"three applied migrations and two reverted": {
			migrations: map[string]string{
				"mig_001": "table_001",
				"mig_002": "table_002",
				"mig_003": "table_003",
			},
			downNumber:        2,
			expHistory:        []string{"mig_001"},
			expExistingTables: []string{"table_001"},
			expDeletedTables:  []string{"table_003", "table_002"},
		},
		"three applied migrations and three reverted": {
			migrations: map[string]string{
				"mig_001": "table_001",
				"mig_002": "table_002",
				"mig_003": "table_003",
			},
			downNumber:        3,
			expHistory:        []string{},
			expExistingTables: []string{},
			expDeletedTables:  []string{"table_003", "table_002", "table_001"},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resetMigrations()
			db = getDB(fmt.Sprintf("tdr_%s", strings.ReplaceAll(testName, " ", "_")))

			for migName, table := range tc.migrations {
				queryUp := fmt.Sprintf("CREATE TABLE %s (id INT)", table)
				queryDown := fmt.Sprintf("DROP TABLE %s", table)
				Add(migName, func(db *sql.DB) error {
					_, err := db.Exec(queryUp)
					return err
				}, func(db *sql.DB) error {
					_, err := db.Exec(queryDown)
					return err
				})
			}

			err := Apply()
			assert.Nil(t, err, "Unexpected error during apply")

			err = Down(tc.downNumber)
			assert.Nil(t, err, "Unexpected error during down")

			history, err := History()
			assert.Nil(t, err, "Unexpected error during history")

			assert.EqualValues(t, tc.expHistory, history)

			for _, tableName := range tc.expExistingTables {
				if !tableExists(db, tableName) {
					t.Errorf("I've expected to see table '%s' but it doesn't exist", tableName)
				}
			}

			for _, tableName := range tc.expDeletedTables {
				if tableExists(db, tableName) {
					t.Errorf("I've expected not to see table '%s' but it exists", tableName)
				}
			}
		})
	}
}
