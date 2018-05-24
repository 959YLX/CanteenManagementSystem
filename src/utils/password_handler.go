package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// EncodePassword 编码password
func EncodePassword(password string) (*string, error) {
	hash := sha256.New()
	if _, err := hash.Write([]byte(password)); err != nil {
		return nil, err
	}
	encodedHexStringPassword := hex.EncodeToString(hash.Sum(nil))
	return &encodedHexStringPassword, nil
}
