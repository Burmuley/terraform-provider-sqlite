package sqlite

import (
	"database/sql"
	"errors"
	"sync"

	_ "modernc.org/sqlite"
)

type sqLiteWrapper struct {
	*sync.Mutex
	db     *sql.DB
	dbPath string
}

func NewSqLiteWrapper() *sqLiteWrapper {
	s := &sqLiteWrapper{
		Mutex: &sync.Mutex{},
	}

	return s
}

func (s *sqLiteWrapper) Open(path string) error {
	var err error
	var db *sql.DB

	// if DB has been already opened - do nothing
	if path == s.dbPath && s.db != nil {
		return nil
	}

	db, err = sql.Open("sqlite", path)

	if err != nil {
		return err
	}

	s.db = db
	s.dbPath = path

	return nil
}

func (s *sqLiteWrapper) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if s.db == nil {
		return nil, errors.New("database not initialized")
	}

	s.Lock()
	defer s.Unlock()
	return s.db.Query(query, args...)
}

func (s *sqLiteWrapper) QueryRow(query string, args ...interface{}) (*sql.Row, error) {
	if s.db == nil {
		return nil, errors.New("database not initialized")
	}

	s.Lock()
	defer s.Unlock()
	row := s.db.QueryRow(query, args...)
	return row, row.Err()
}

func (s *sqLiteWrapper) Exec(query string, args ...interface{}) (sql.Result, error) {
	if s.db == nil {
		return nil, errors.New("database not initialized")
	}

	s.Lock()
	defer s.Unlock()
	return s.db.Exec(query, args...)
}

func (s *sqLiteWrapper) Ping() error {
	return s.db.Ping()
}
