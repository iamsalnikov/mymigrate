// +build integration

package mymigrate

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHistory_CreatesMigrationTable(t *testing.T) {
	resetMigrations()
	db = getDB("test_history_creates_migration_table")
	_, _ = History()
	row := db.QueryRow("SHOW TABLES LIKE '%mymigrations%'")
	var s string
	row.Scan(&s)
	fmt.Println(s)
	if s != "mymigrations" {
		t.Errorf("Не удалось найти таблицу mymigrations после вызова Apply() на чистой БД")
	}
	db = nil
}

func TestHistory_OrderList(t *testing.T) {
	resetMigrations()
	db = getDB("test_history_order_list")
	_ = defaultMarkAppliedFunc(db, "hello 1")
	time.Sleep(1 * time.Second)
	_ = defaultMarkAppliedFunc(db, "hello 2")
	time.Sleep(1 * time.Second)
	_ = defaultMarkAppliedFunc(db, "hello 3")

	list, _ := History()
	assert.EqualValues(t, []string{"hello 3", "hello 2", "hello 1"}, list)
}
