package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

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
