package share

import (
	"fmt"
	"os"
)

// LoadTestPNG loads a PNG file from the specified path
func LoadTestPNG(filepath string) ([]byte, error) {
	// Try the filepath as-is first
	data, err := os.ReadFile(filepath)
	if err == nil {
		return data, nil
	}

	// If that fails, try relative to the server directory
	serverPath := "../" + filepath
	data, err = os.ReadFile(serverPath)
	if err == nil {
		return data, nil
	}

	return nil, fmt.Errorf("failed to load PNG file from %s or %s: %v", filepath, serverPath, err)
}
