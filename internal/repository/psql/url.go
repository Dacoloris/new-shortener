package psql

import (
	"database/sql"

	_ "github.com/jackc/pgx"
)

func PingDB(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return err
	}

	return nil
}
