package utils

import (
	"time"

	"golang.org/x/exp/rand"
)

func generateUniqueID() string {

	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 25 // Adjust the length as needed

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}

	return string(b)
}
