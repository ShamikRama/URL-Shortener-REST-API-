package psql

import (
	"Rest/internal/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(cfg config.Config) (*Storage, error) {
	const op = "storage.psql.New"

	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	//defer db.Close()

	// Создание таблицы
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL
		);
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	// Создание индекса
	createIndexSQL := `
		CREATE INDEX IF NOT EXISTS idx_alias ON urls(alias);
	`
	_, err = db.Exec(createIndexSQL)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}
