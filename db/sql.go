package db

import (
	"github.com/jmoiron/sqlx"
)

type DB struct {
	db *sqlx.DB
}

func NewDB(db *sqlx.DB) *DB {
	return &DB{
		db: db,
	}
}

func (db *DB) Close() error {
	return db.db.Close()
}
