package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `
	CREATE TABLE scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(256) NOT NULL DEFAULT '',
		comment TEXT NOT NULL DEFAULT "",
		repeat VARCHAR(256) NOT NULL DEFAULT ''
	);
	CREATE INDEX task_date ON scheduler (date);`

func Init(dbFile string) error {
	var err error

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("Failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("Failed to ping database: %w", err)
	}

	var tableExists bool
	err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='scheduler'").Scan(&tableExists)
	if err != nil || !tableExists {
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("Failed to create tables: %w", err)
		}
	}

	return nil
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
