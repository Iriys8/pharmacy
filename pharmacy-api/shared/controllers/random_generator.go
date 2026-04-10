package controllers

import (
	"crypto/rand"
	"encoding/hex"
)

func Random_gen() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
