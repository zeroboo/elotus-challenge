package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

var DB *sql.DB

// InitDB initializes the SQLite database connection and creates tables
func InitDB() error {
	// Get database path from environment or use default
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./challenge.db"
	}

	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		return err
	}

	// Create tables
	if err = createTables(); err != nil {
		return err
	}

	log.Info().Str("db_path", dbPath).Msg("Database initialized successfully")
	return nil
}

// createTables creates the required tables if they don't exist
func createTables() error {
	// Users table
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(50) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Execute table creation
	tables := []string{userTable, fileTable, tokenTable}
	for _, table := range tables {
		if _, err := DB.Exec(table); err != nil {
			return err
		}
	}

	log.Info().Msg("Database tables created successfully")
	return nil
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Info().Msg("Database connection closed")
	}
}
