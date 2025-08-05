package repository

import (
	"database/sql"
	"elotuschallenge/database"
	"elotuschallenge/models"
)

type SQLiteUserRepository struct{}

func NewSQLiteUserRepository() IUser {
	return &SQLiteUserRepository{}
}

// CreateUser inserts a new user into the database and returns the user with ID
func (r *SQLiteUserRepository) CreateUser(user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (username, password_hash, created_at) 
		VALUES (?, ?, CURRENT_TIMESTAMP)
	`

	result, err := database.DB.Exec(query, user.Username, user.PasswordHash)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user.ID = int(id)
	return user, nil
}

// UserExists checks if a user with the given username already exists
func (r *SQLiteUserRepository) UserExists(username string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE username = ?"
	var count int
	err := database.DB.QueryRow(query, username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetUserByUsername retrieves a user by username
func (r *SQLiteUserRepository) GetUserByUsername(username string) (*models.User, error) {
	query := "SELECT id, username, password_hash FROM users WHERE username = ?"
	var user models.User
	err := database.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}
	return &user, nil
}
