package pggood

import "github.com/jmoiron/sqlx"

type PgGoodStorage struct {
	db *sqlx.DB
}

func NewPgGoodStorage(db *sqlx.DB) *PgGoodStorage {
	return &PgGoodStorage{
		db: db,
	}
}
