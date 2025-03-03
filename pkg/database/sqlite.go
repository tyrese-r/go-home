package database

import (
	"database/sql"
	_ "github.com/glebarez/sqlite"
)

// NewSQLiteDB creates and initializes a new SQLite database connection
func NewSQLiteDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Initialize database schema
	if err := initSchema(db); err != nil {
		return nil, err
	}

	return db, nil
}

// initSchema creates necessary tables if they don't exist
func initSchema(db *sql.DB) error {
	// Create devices table
	devicesTableDDL := `
	CREATE TABLE IF NOT EXISTS devices (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		device_type TEXT NOT NULL,
		owned_by TEXT NOT NULL,
		is_online BOOLEAN DEFAULT FALSE,
		last_alarm_reason TEXT,
		last_alarm_time TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(devicesTableDDL); err != nil {
		return err
	}

	return nil
}
