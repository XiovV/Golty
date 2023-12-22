package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const DefaultQueryTimeout = 5

type Repository struct {
	db *sqlx.DB
}

func New(dbPath string) (*Repository, error) {
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	db.MustExec("PRAGMA foreign_keys = ON;")

	return &Repository{db: db}, nil
}

func newBackgroundContext(duration int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
}

type BoolAsInt bool

func (b *BoolAsInt) Scan(value any) error {
	i, ok := value.(int64)
	if !ok {
		return fmt.Errorf("expected int64, got %T", value)
	}

	*b = i != 0
	return nil
}
