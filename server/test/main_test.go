// test/main_test.go
package test

import (
	"os"
	"testing"

	"elotuschallenge/database"
	"elotuschallenge/internal"
)

func TestMain(m *testing.M) {
	// Setup before all tests
	setup()

	// Run tests
	code := m.Run()

	// Cleanup after all tests
	cleanup()

	os.Exit(code)
}

func setup() {
	// Setup: Use a test database
	os.Setenv("DB_PATH", ":memory:")

	// Initialize test database
	if err := database.InitDB(); err != nil {
		panic("Failed to initialize test database: " + err.Error())
	}
	internal.InitServices()
}

func cleanup() {
	// Cleanup
	database.CloseDB()
}
