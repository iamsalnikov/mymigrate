package mymigrate

import (
	"context"
	"database/sql"
	"time"
)

// UpFunc is a function that ups migration
type UpFunc func(db *sql.DB) error

// DownFunc is a function that downs migration
type DownFunc func(db *sql.DB) error

type mig struct {
	name string
	up   UpFunc
	down DownFunc
}

// DbProvider - interface for interacting with the database
type DbProvider interface {
	GetDb() *sql.DB
	CreateMigrationsTable() error
	GetApplied(context.Context) ([]string, error)
	MarkApplied(context.Context, string, time.Time) error
	DeleteApplied(context.Context, string) error
}
