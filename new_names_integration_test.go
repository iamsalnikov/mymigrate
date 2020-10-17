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
	if !tableExists(db, "mymigrations") {
		t.Errorf("Не удалось найти таблицу mymigrations после вызова NewNames() на чистой БД")
	}
	db = nil
}

func TestNewNames(t *testing.T) {
	type testCase struct {
		AppliedNames []string
		ToAdd        []string
		Exp          []string
	}

	testCases := []testCase{
		{
			AppliedNames: []string{},
			ToAdd: []string{
				"m_001",
				"m_002",
				"m_003",
			},
			Exp: []string{"m_001", "m_002", "m_003"},
		},
		{
			AppliedNames: []string{"m_005"},
			ToAdd: []string{
				"m_001",
				"m_002",
				"m_003",
			},
			Exp: []string{"m_001", "m_002", "m_003"},
		},
		{
			AppliedNames: []string{"m_002"},
			ToAdd: []string{
				"m_001",
				"m_002",
				"m_003",
			},
			Exp: []string{"m_001", "m_003"},
		},
		{
			AppliedNames: []string{"m_002", "m_001", "m_003"},
			ToAdd: []string{
				"m_001",
				"m_002",
				"m_003",
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

			for _, name := range tc.ToAdd {
				Add(
					name,
					func(db *sql.DB) error { return nil },
					func(db *sql.DB) error { return nil },
				)
			}

			newNames, _ := NewNames()
			assert.EqualValues(t, tc.Exp, newNames)
		})
	}

}
