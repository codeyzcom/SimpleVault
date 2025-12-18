package vault

import (
	"crypto/rand"
	"crypto/sha256"
)

func sha(b []byte) []byte {
	h := sha256.Sum256(b)
	return h[:]
}

func GeneratePassword(length int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"
	if length < 6 {
		length = 6
	}

	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	for i := range b {
		b[i] = chars[int(b[i])%len(chars)]
	}
	return string(b), nil
}
