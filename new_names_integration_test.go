// +build integration

package mymigrate

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNames_CreatesMigrationTable(t *testing.T) {
	resetMigrations()
	db = getDB("test_new_names_creates_migration_table")
	_, _ = NewNames()
	row := db.QueryRow("SHOW TABLES LIKE '%mymigrations%'")
	var s string
	row.Scan(&s)
	fmt.Println(s)
	if s != "mymigrations" {
		t.Errorf("Не удалось найти таблицу mymigrations после вызова NewNames() на чистой БД")
	}
	db = nil
}

func TestNewNames(t *testing.T) {
	type testCase struct {
		AppliedNames []string
		ToAdd        map[string]UpFunc
		Exp          []string
	}

	testCases := []testCase{
		{
			AppliedNames: []string{},
			ToAdd: map[string]UpFunc{
				"m_001": func(db *sql.DB) error {
					return nil
				},
				"m_002": func(db *sql.DB) error {
					return nil
				},
				"m_003": func(db *sql.DB) error {
					return nil
				},
			},
			Exp: []string{"m_001", "m_002", "m_003"},
		},
		{
			AppliedNames: []string{"m_005"},
			ToAdd: map[string]UpFunc{
				"m_001": func(db *sql.DB) error {
					return nil
				},
				"m_002": func(db *sql.DB) error {
					return nil
				},
				"m_003": func(db *sql.DB) error {
					return nil
				},
			},
			Exp: []string{"m_001", "m_002", "m_003"},
		},
		{
			AppliedNames: []string{"m_002"},
			ToAdd: map[string]UpFunc{
				"m_001": func(db *sql.DB) error {
					return nil
				},
				"m_002": func(db *sql.DB) error {
					return nil
				},
				"m_003": func(db *sql.DB) error {
					return nil
				},
			},
			Exp: []string{"m_001", "m_003"},
		},
		{
			AppliedNames: []string{"m_002", "m_001", "m_003"},
			ToAdd: map[string]UpFunc{
				"m_001": func(db *sql.DB) error {
					return nil
				},
				"m_002": func(db *sql.DB) error {
					return nil
				},
				"m_003": func(db *sql.DB) error {
					return nil
				},
			},
			Exp: []string{},
		},
	}

	for i, tc := range testCases {
		resetMigrations()
		dbn := fmt.Sprintf("test_new_names_case_%d", i)
		t.Run(dbn, func(t *testing.T) {
			db = getDB(dbn)

			for _, appliedName := range tc.AppliedNames {
				_ = defaultMarkAppliedFunc(db, appliedName)
			}

			for name, up := range tc.ToAdd {
				Add(name, up)
			}

			newNames, _ := NewNames()
			assert.EqualValues(t, tc.Exp, newNames)
		})
	}

}
