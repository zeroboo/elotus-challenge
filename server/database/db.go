package database

import (
	"database/sql"
	"os"

	"github.com/rs/zerolog/log"
	_ "modernc.org/sqlite"
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
	DB, err = sql.Open("sqlite", dbPath)
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

	// Files table
	fileTable := `
	CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filename VARCHAR(255) NOT NULL,
		original_name VARCHAR(255) NOT NULL,
		content_type VARCHAR(100) NOT NULL,
		size INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		upload_path VARCHAR(500) NOT NULL,
		user_agent VARCHAR(500),
		ip_address VARCHAR(45),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	// Optional: Token blacklist for revocation
	tokenTable := `
	CREATE TABLE IF NOT EXISTS revoked_tokens (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		token_hash VARCHAR(255) UNIQUE NOT NULL,
		revoked_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		expires_at DATETIME NOT NULL
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
