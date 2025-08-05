package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const RandomStringCharset = "ABCDEFGHIJKLMN123456789"

// GenerateRandomString generates a random string of specified length from a predefined character set
func GenerateRandomString(length int) string {
	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(RandomStringCharset)))

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			// Fallback to a simple format if random generation fails
			return fmt.Sprintf("ERR-%d", i)
		}
		result[i] = RandomStringCharset[randomIndex.Int64()]
	}

	return string(result)
}
