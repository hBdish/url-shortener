package postgres

import (
	"database/sql"
	"fmt"
	"url-shortener/internal/config"
)

type Storage struct {
	db *sql.DB
}

func New(dbCnf config.Db) (*Storage, error) {
	const fn = "storage.postgres.New"

	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbCnf.Host, dbCnf.Port, dbCnf.User, dbCnf.Password, dbCnf.Dbname)

	db, err := sql.Open("postgres", dbInfo)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
    	id SERIAL PRIMARY KEY ,
    	alias TEXT NOT NULL UNIQUE,
    	url TEXT NOT NULL
		);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const fn = "storage.postgres.SaveURL"

	var id int64
	err := s.db.QueryRow("INSERT INTO url(url, alias) VALUES ($1, $2) RETURNING id", urlToSave, alias).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return id, err
}

func (s *Storage) GetURL(alias string) (string, error) {
	const fn = "storage.postgres.GetURL"

	var url string
	err := s.db.QueryRow("SELECT url FROM url WHERE alias = $1", alias).Scan(&url)

	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	return url, err
}
