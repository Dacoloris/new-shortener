package psql

import (
	"context"
	"database/sql"
	"errors"
	"new-shortner/internal/domain"
)

var (
	ErrNotFound = errors.New("not found")
)

type Storage struct {
	conn *sql.DB
	dsn  string
}

func New(db *sql.DB, dsn string) *Storage {
	return &Storage{
		conn: db,
		dsn:  dsn,
	}
}

func (s *Storage) Create(ctx context.Context, url domain.URL) error {
	query := `INSERT INTO urls (user_id, original_url, short_url)
				  VALUES ($1, $2, $3)`

	_, err := s.conn.ExecContext(ctx, query, url.UserID, url.Original, url.Short)

	return err
}

func (s *Storage) GetOriginalByShort(ctx context.Context, shortURL string) (string, error) {
	query := `SELECT original_url FROM urls WHERE short_url=$1 LIMIT 1`
	row := s.conn.QueryRowContext(ctx, query, shortURL)
	var original string
	row.Scan(&original)
	if original == "" {
		return "", ErrNotFound
	}

	return original, nil
}

func (s *Storage) GetAllURLsByUserID(ctx context.Context, userID string) ([]domain.URL, error) {
	res := make([]domain.URL, 0)

	query := `SELECT original_url, short_url FROM urls WHERE user_id=$1`
	rows, err := s.conn.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u domain.URL
		err = rows.Scan(&u.Original, &u.Short)
		if err != nil {
			return nil, err
		}

		res = append(res, u)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return res, nil
}

func (s *Storage) CreateBatch(ctx context.Context, urls []domain.URL) error {
	tx, err := s.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO urls (user_id, original_url, short_url)
				  VALUES ($1, $2, $3)`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	for _, url := range urls {
		_, err = stmt.ExecContext(ctx, url.UserID, url.Original, url.Short)
		if err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}
