package helper

import (
	"crypto/rand"
	"math/big"
	"strings"
)

func GenerateOTPCode(length int) (string, error) {
	seed := "012345679"
	byteSlice := make([]byte, length)

	for i := 0; i < length; i++ {
		max := big.NewInt(int64(len(seed)))
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}

		byteSlice[i] = seed[num.Int64()]
	}

	return string(byteSlice), nil
}

func GenerateREFCode(length int) (string, error) {
	seed := "0123456789abcdefghijklmnopqrstuvwxyz"
	byteSlice := make([]byte, length)

	for i := 0; i < length; i++ {
		max := big.NewInt(int64(len(seed)))
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}

		byteSlice[i] = seed[num.Int64()]
	}

	return strings.ToUpper(string(byteSlice)), nil
}

func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email // Return original if not a valid email format
	}

	username := parts[0]
	domain := parts[1]

	if len(username) == 0 {
		return email
	}

	if len(username) == 1 {
		return username + "@" + domain
	}

	// Create masked username: first char + x's for the rest
	maskedUsername := string(username[0]) + strings.Repeat("x", len(username)-1)

	return maskedUsername + "@" + domain
}
