package postgres_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/iamsalnikov/mymigrate/provider"
	"github.com/iamsalnikov/mymigrate/provider/postgres"
	"github.com/stretchr/testify/assert"
)

func TestPsqlProvider_CreateMigrationsTable(t *testing.T) {

	cases := map[string]struct {
		execError   error
		expectQuery string
		expectErr   error
	}{
		"All is ok": {
			execError:   nil,
			expectQuery: fmt.Sprintf("create table if not exists %s ( name varchar(500) not null constraint %s_pk primary key, time timestamp ); create unique index if not exists %s_name_uindex on %s (name);", provider.DefaultTableName, provider.DefaultTableName, provider.DefaultTableName, provider.DefaultTableName),
			expectErr:   nil,
		},
		"db error": {
			execError:   errors.New("some db error"),
			expectQuery: fmt.Sprintf("create table if not exists %s ( name varchar(500) not null constraint %s_pk primary key, time timestamp ); create unique index if not exists %s_name_uindex on %s (name);", provider.DefaultTableName, provider.DefaultTableName, provider.DefaultTableName, provider.DefaultTableName),
			expectErr:   errors.New("some db error"),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			defer db.Close()
			assert.NoError(t, err)

			mock.ExpectExec(c.expectQuery).
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(c.execError)

			p := postgres.NewPsqlProvider(db)
			err = p.CreateMigrationsTable()

			assert.Equal(t, c.expectErr, err)
		})
	}

}

func TestPsqlProvider_GetDb(t *testing.T) {
	db, _, err := sqlmock.New()
	defer db.Close()
	assert.NoError(t, err)

	p := postgres.NewPsqlProvider(db)

	assert.Equal(t, db, p.GetDb())
}

func TestPsqlProvider_GetApplied(t *testing.T) {
	cases := map[string]struct {
		execError error
		execRows  *sqlmock.Rows

		expectQuery  string
		expectErr    error
		expectResult []string
	}{
		"empty migration table": {
			execError:    nil,
			execRows:     sqlmock.NewRows([]string{"name"}),
			expectQuery:  fmt.Sprintf("SELECT name FROM %s ORDER BY time DESC, name DESC", provider.DefaultTableName),
			expectErr:    nil,
			expectResult: []string{},
		},
		"all is ok": {
			execError:    nil,
			execRows:     sqlmock.NewRows([]string{"name"}).AddRow("migration_1").AddRow("migration_2"),
			expectQuery:  fmt.Sprintf("SELECT name FROM %s ORDER BY time DESC, name DESC", provider.DefaultTableName),
			expectErr:    nil,
			expectResult: []string{"migration_1", "migration_2"},
		},
		"db error": {
			execError:    errors.New("some db error"),
			execRows:     sqlmock.NewRows([]string{"name"}),
			expectQuery:  fmt.Sprintf("SELECT name FROM %s ORDER BY time DESC, name DESC", provider.DefaultTableName),
			expectErr:    errors.New("some db error"),
			expectResult: nil,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			defer db.Close()
			assert.NoError(t, err)

			mock.ExpectQuery(c.expectQuery).WillReturnRows(c.execRows).WillReturnError(c.execError)

			p := postgres.NewPsqlProvider(db)

			res, err := p.GetApplied(context.Background())

			assert.Equal(t, c.expectErr, err)
			assert.Equal(t, c.expectResult, res)
		})
	}
}

func TestPsqlProvider_MarkApplied(t *testing.T) {
	now := time.Now()
	cases := map[string]struct {
		name string
		time time.Time

		execError error

		expectQuery string
		expectErr   error
		expectArgs  []interface{}
	}{
		"all is ok": {
			name:        "migration_1",
			time:        now,
			execError:   nil,
			expectQuery: fmt.Sprintf("INSERT INTO %s (name, time) VALUES ($1, $2)", provider.DefaultTableName),
			expectArgs:  []interface{}{"migration_1", now},
			expectErr:   nil,
		},
		"db error": {
			name:        "migration_2",
			time:        now,
			execError:   errors.New("some db error"),
			expectQuery: fmt.Sprintf("INSERT INTO %s (name, time) VALUES ($1, $2)", provider.DefaultTableName),
			expectArgs:  []interface{}{"migration_2", now},
			expectErr:   errors.New("some db error"),
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			defer db.Close()
			assert.NoError(t, err)

			mock.ExpectExec(c.expectQuery).
				WithArgs(c.expectArgs[0], c.expectArgs[1]).
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(c.execError)

			p := postgres.NewPsqlProvider(db)

			err = p.MarkApplied(context.Background(), c.name, c.time)

			assert.Equal(t, c.expectErr, err)
		})
	}
}

func TestPsqlProvider_DeleteApplied(t *testing.T) {
	cases := map[string]struct {
		name string
		time time.Time

		execError error

		expectQuery string
		expectErr   error
		expectArgs  []interface{}
	}{
		"all is ok": {
			name:        "migration_1",
			execError:   nil,
			expectQuery: fmt.Sprintf("DELETE FROM %s WHERE name=$1", provider.DefaultTableName),
			expectArgs:  []interface{}{"migration_1"},
			expectErr:   nil,
		},
		"db error": {
			name:        "migration_2",
			execError:   errors.New("some db error"),
			expectQuery: fmt.Sprintf("DELETE FROM %s WHERE name=$1", provider.DefaultTableName),
			expectArgs:  []interface{}{"migration_2"},
			expectErr:   errors.New("some db error"),
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			defer db.Close()
			assert.NoError(t, err)

			mock.ExpectExec(c.expectQuery).
				WithArgs(c.expectArgs[0]).
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(c.execError)

			p := postgres.NewPsqlProvider(db)

			err = p.DeleteApplied(context.Background(), c.name)

			assert.Equal(t, c.expectErr, err)
		})
	}
}
