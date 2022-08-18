package repository

import "github.com/jmoiron/sqlx"

type ConnManager interface {
	GetReader() *sqlx.DB
	GetWriter() *sqlx.DB
	CloseAll()
}
