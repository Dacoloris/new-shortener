package initdb

import (
	"database/sql"
)

func InitDB(db *sql.DB) error {
	sqlCreateDB := `CREATE TABLE IF NOT EXISTS urls (
								id serial PRIMARY KEY,
								user_id uuid NOT NULL ,	
								original_url VARCHAR NOT NULL, 
								short_url VARCHAR NOT NULL UNIQUE
					);`
	_, err := db.Exec(sqlCreateDB)
	if err != nil {
		return err
	}

	return nil
}
