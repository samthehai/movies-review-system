package repository_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/samthehai/ml-backend-test-samthehai/internal/movie/interfaceadapters/repository"
)

var favoritesTableRows []string = []string{"user_id", "movie_id", "created_at", "updated_at"}
var moviesTableRows []string = []string{"id", "original_title", "original_language", "overview",
	"poster_path", "backdrop_path", "adult", "release_date", "budget", "revenue", "created_at", "updated_at"}

type connManager struct {
	db *sqlx.DB
}

func (c connManager) GetReader() *sqlx.DB {
	return c.db
}

func (c connManager) GetWriter() *sqlx.DB {
	return c.db
}

func (c connManager) CloseAll() {
	if c.db == nil {
		return
	}

	c.db.Close()
}

func newConnManager() (repository.ConnManager, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	return connManager{sqlx.NewDb(db, "mysql")}, mock, nil
}

func initMockConnManager(t *testing.T, mocks func(mock sqlmock.Sqlmock)) (repository.ConnManager, func()) {
	t.Helper()
	manager, mock, err := newConnManager()
	if err != nil {
		t.Errorf("failed to init database connection manager: %v", err)
	}

	mocks(mock)

	return manager, func() {
		t.Helper()
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfullfilled expectations: %v", err)
		}
		manager.CloseAll()
	}
}
