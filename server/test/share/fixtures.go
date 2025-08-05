package share

import (
	"elotuschallenge/models"
	"elotuschallenge/services"
	"fmt"
	"time"
)

// CreateTestUser creates a test user for testing purposes
func CreateTestUser(id int, username string) *models.User {
	return &models.User{
		ID:           id,
		Username:     username,
		PasswordHash: "$2a$10$examplehash", // Example bcrypt hash
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// CreateTestUsers creates multiple test users
func CreateTestUsers(count int) []*models.User {
	users := make([]*models.User, count)
	for i := 0; i < count; i++ {
		users[i] = CreateTestUser(i+1, fmt.Sprintf("testuser%d", i+1))
	}
	return users
}

// CreateTestJWTService creates a JWT service for testing
func CreateTestJWTService() services.ITokenManager {
	return services.NewTokenManager("test-secret-key", 24*60*60) // 24 hours expiration
}
