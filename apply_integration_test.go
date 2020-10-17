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
	_, _ = Apply()
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
		appliedNames        []string
		toAdd               []string
		expAppliedToHistory []string
		expJustApplied      []string
	}

	testCases := []testCase{
		{
			appliedNames:        []string{},
			toAdd:               []string{"m_001"},
			expAppliedToHistory: []string{"m_001"},
			expJustApplied:      []string{"m_001"},
		},
		{
			appliedNames:        []string{"m_001", "m_002"},
			toAdd:               []string{},
			expAppliedToHistory: []string{"m_002", "m_001"},
			expJustApplied:      []string{},
		},
		{
			appliedNames:        []string{"m_001", "m_002"},
			toAdd:               []string{"m_001", "m_002"},
			expAppliedToHistory: []string{"m_002", "m_001"},
			expJustApplied:      []string{},
		},
		{
			appliedNames:        []string{"m_001", "m_002"},
			toAdd:               []string{"m_001", "m_002", "m_003", "m_004"},
			expAppliedToHistory: []string{"m_004", "m_003", "m_002", "m_001"},
			expJustApplied:      []string{"m_003", "m_004"},
		},
	}

	for i, tc := range testCases {
		resetMigrations()
		dbn := fmt.Sprintf("test_apply_case_%d_%d", i, time.Now().UnixNano())
		t.Run(dbn, func(t *testing.T) {
			db = getDB(dbn)

			for _, appliedName := range tc.appliedNames {
				_ = defaultMarkAppliedFunc(db, appliedName)
			}

			for _, name := range tc.toAdd {
				Add(name, func(db *sql.DB) error {
					return nil
				}, func(db *sql.DB) error {
					return nil
				})
			}

			justApplied, err := Apply()
			assert.Nil(t, err, "Unexpected error during apply")
			applied, _ := History()

			assert.EqualValues(t, tc.expAppliedToHistory, applied)
			assert.EqualValues(t, tc.expJustApplied, justApplied)
		})
	}
}
