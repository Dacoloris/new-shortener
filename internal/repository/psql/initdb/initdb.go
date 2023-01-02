package initdb

import (
	"database/sql"
)

func InitDB(db *sql.DB) error {
	sqlCreateDB := `CREATE TABLE IF NOT EXISTS urls (
								id serial PRIMARY KEY, 	
								origin_url VARCHAR NOT NULL, 
								short_url VARCHAR NOT NULL UNIQUE
					);`
	_, err := db.Exec(sqlCreateDB)
	if err != nil {
		return err
	}

	return nil
}
