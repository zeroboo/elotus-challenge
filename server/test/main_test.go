// test/main_test.go
package test

import (
	"os"
	"testing"
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
}

func cleanup() {
}
