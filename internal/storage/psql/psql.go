package psql

import (
	"Rest/internal/config"
	"Rest/internal/storage"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

var ErrConstraintCode = 23505

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

func (s *Storage) SaveUrl(urlToSave string, alias string) (int64, error) {
	const op = "storage.psql.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO urls(url, alias) VALUES($1, $2)")
	if err != nil {
		return 0, fmt.Errorf("%s : %w", op, err)
	}

	var id int64
	err = stmt.QueryRow(urlToSave, alias).Scan(&id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "ErrConstraintCode" {
			return 0, fmt.Errorf("%s : %w", op, storage.ErrUrlExist)
		}
		return 0, fmt.Errorf("%s : %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	const op = "storage.psql.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM urls WHERE alias = $1 )")
	if err != nil {
		return "", fmt.Errorf("%s : %w", op, err)
	}

	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrUrlNotFound
		}
		return "", fmt.Errorf("%s : %w", op, err)
	}

	return resURL, nil
}

func (s *Storage) DeleteUrl(alias string) error {
	const op = "storage.psql.DeleteURL"
	stmt, err := s.db.Prepare("DELETE FROM urls WHERE alias = $1")
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}

	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}

	return nil
}
