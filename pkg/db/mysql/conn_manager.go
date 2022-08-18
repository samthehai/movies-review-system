package mysql

import (
	"fmt"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/samthehai/ml-backend-test-samthehai/config"
)

const (
	maxOpenConns    = 60
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

type connManager struct {
	writer  *sqlx.DB
	readers []*sqlx.DB
}

func NewConnManager(cfg *config.Config) (*connManager, func(), error) {
	readerDBs := make([]*sqlx.DB, 0, len(cfg.MySQL.ReaderDataSources))
	clean := func() {
		for _, c := range readerDBs {
			_ = c.Close()
		}
	}

	for _, ds := range cfg.MySQL.ReaderDataSources {
		db, err := sqlx.Connect("mysql", fmt.Sprintf("%s?parseTime=true", ds))
		if err != nil {
			clean()
			return nil, nil, fmt.Errorf("init db connection: %w", err)
		}

		db.SetMaxOpenConns(maxOpenConns)
		db.SetConnMaxLifetime(connMaxLifetime * time.Second)
		db.SetMaxIdleConns(maxIdleConns)
		db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)
		if err = db.Ping(); err != nil {
			clean()
			return nil, nil, fmt.Errorf("verify db connection: %w", err)
		}

		readerDBs = append(readerDBs, db)
	}

	writerDB, err := sqlx.Connect("mysql", fmt.Sprintf("%s?parseTime=true", cfg.MySQL.WriterDataSource))
	if err != nil {
		clean()
		return nil, nil, fmt.Errorf("init db connection: %w", err)
	}

	writerDB.SetMaxOpenConns(maxOpenConns)
	writerDB.SetConnMaxLifetime(connMaxLifetime * time.Second)
	writerDB.SetMaxIdleConns(maxIdleConns)
	writerDB.SetConnMaxIdleTime(connMaxIdleTime * time.Second)
	if err = writerDB.Ping(); err != nil {
		clean()
		return nil, nil, fmt.Errorf("verify db connection: %w", err)
	}

	connMng := &connManager{
		readers: readerDBs,
		writer:  writerDB,
	}

	return connMng, func() { connMng.CloseAll() }, nil
}

func (m *connManager) GetReader() *sqlx.DB {
	rand.Seed(time.Now().UnixNano())
	return m.readers[rand.Intn(len(m.readers))]
}

func (m *connManager) GetWriter() *sqlx.DB {
	return m.writer
}

func (m *connManager) CloseAll() {
	var c *sqlx.DB
	for _, c = range m.readers {
		_ = c.Close()
	}

	if err := m.writer.Close(); err != nil {
		panic(err)
	}
}
