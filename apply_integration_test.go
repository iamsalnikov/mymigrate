// +build integration

package mymigrate

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestApply_CreatesMigrationTable(t *testing.T) {
	resetMigrations()
	db = getDB("test_apply_creates_migration_table")
	_ = Apply()
	row := db.QueryRow("SHOW TABLES LIKE '%mymigrations%'")
	var s string
	row.Scan(&s)
	fmt.Println(s)
	if s != "mymigrations" {
		t.Errorf("Не удалось найти таблицу mymigrations после вызова Apply() на чистой БД")
	}
	db = nil
}

func TestApply(t *testing.T) {
	type testCase struct {
		AppliedNames []string
		ToAdd        []string
		Exp          []string
	}

	testCases := []testCase{
		{
			AppliedNames: []string{},
			ToAdd:        []string{"m_001"},
			Exp:          []string{"m_001"},
		},
		{
			AppliedNames: []string{"m_001", "m_002"},
			ToAdd:        []string{},
			Exp:          []string{"m_002", "m_001"},
		},
		{
			AppliedNames: []string{"m_001", "m_002"},
			ToAdd:        []string{"m_001", "m_002"},
			Exp:          []string{"m_002", "m_001"},
		},
		{
			AppliedNames: []string{"m_001", "m_002"},
			ToAdd:        []string{"m_001", "m_002", "m_003", "m_004"},
			Exp:          []string{"m_004", "m_003", "m_002", "m_001"},
		},
	}

	for i, tc := range testCases {
		resetMigrations()
		dbn := fmt.Sprintf("test_apply_case_%d_%d", i, time.Now().UnixNano())
		t.Run(dbn, func(t *testing.T) {
			db = getDB(dbn)

			for _, appliedName := range tc.AppliedNames {
				_ = defaultMarkAppliedFunc(db, appliedName)
			}

			for _, name := range tc.ToAdd {
				Add(name, func(db *sql.DB) error {
					return nil
				}, func(db *sql.DB) error {
					return nil
				})
			}

			_ = Apply()
			applied, _ := History()

			assert.EqualValues(t, tc.Exp, applied)
		})
	}
}
