package ping

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx"
)

func Ping(ctx context.Context, dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err = db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}
