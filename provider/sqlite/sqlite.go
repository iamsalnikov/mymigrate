package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/iamsalnikov/mymigrate/provider"
	"time"
)

// Provider - migration provider for sqlite db
type Provider struct {
	db *sql.DB
}

// NewSqliteProvider - constructor for sqlite Provider
func NewSqliteProvider(db *sql.DB) *Provider {
	return &Provider{db: db}
}

// GetDb - function returning internal db object
func (p *Provider) GetDb() *sql.DB {
	return p.db
}

// CreateMigrationsTable - function creating migration table in db
func (p *Provider) CreateMigrationsTable() error {
	query := fmt.Sprintf(`create table if not exists %s
			(
				name varchar(500) not null constraint table_name_pk primary key,
				time timestamp
			);
		create unique index if not exists %s_name_uindex on %s (name);`, provider.DefaultTableName, provider.DefaultTableName, provider.DefaultTableName)

	_, err := p.db.Exec(query)
	return err
}

// GetApplied - function returning list applied migrations
func (p *Provider) GetApplied(ctx context.Context) ([]string, error) {
	query := fmt.Sprintf("SELECT name FROM %s ORDER BY time DESC, name DESC", provider.DefaultTableName)
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	res := make([]string, 0)
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}

		res = append(res, name)
	}
	return res, nil
}

// MarkApplied - function for mark migration applied
func (p *Provider) MarkApplied(ctx context.Context, name string, t time.Time) error {
	query := fmt.Sprintf("INSERT INTO %s (name, time) VALUES (?, ?)", provider.DefaultTableName)
	_, err := p.db.ExecContext(ctx, query, name, t)
	return err
}

// DeleteApplied - function for delete migration from applied list
func (p *Provider) DeleteApplied(ctx context.Context, name string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE name=?", provider.DefaultTableName)
	_, err := p.db.ExecContext(ctx, query, name)
	return err
}
