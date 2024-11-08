package postgres

import (
	"database/sql"
	"fmt"
	"url-shortener/internal/config"

	_ "github.com/lib/pq"
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
    	id INTEGER PRIMARY KEY ,
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
