package utils

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"golang.org/x/exp/rand"
)

// GenerateUniqueID generates a unique ID based on the input
func GenerateUniqueID(in ...string) string {

	if len(in) == 0 {
		return randomString()
	}

	return generateShortHash(in[0])
}

// generateShortHash generates a unique hash (MD5) for the input string and returns the first 10 characters
func generateShortHash(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])[:10]
}

// randomString generates a random string {
func randomString() string {
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	length := 10 // Adjust the length as needed

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}

	return string(b)
}
