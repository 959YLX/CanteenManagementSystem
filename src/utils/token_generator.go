package utils

import (
	"github.com/satori/go.uuid"
)

// GenerateToken 生成随机token
func GenerateToken() (*string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	stringUUID := u.String()
	return &stringUUID, nil
}
