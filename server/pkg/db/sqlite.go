package db

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

const DefaultQueryTimeout = 5

type DB struct {
	db     *sqlx.DB
	logger *zap.SugaredLogger
}

func New(dbPath string, logger *zap.SugaredLogger) (*DB, error) {
	sqliteDb, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	sqliteDb.MustExec("PRAGMA foreign_keys = ON;")

	db := &DB{db: sqliteDb, logger: logger}

	db.initSchema()

	return db, nil
}

func newBackgroundContext(duration int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
}
