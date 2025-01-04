package utils

import (
	"strconv"
	"time"

	"github.com/mitchellh/hashstructure/v2"
	"golang.org/x/exp/rand"
)

func GenerateUniqueID(in ...string) string {

	if len(in) == 0 {
		return randomString()
	}

	hash, err := hashstructure.Hash(in[0], hashstructure.FormatV2, nil)
	if err != nil {
		return randomString()
	}
	return strconv.FormatUint(hash, 10)
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
