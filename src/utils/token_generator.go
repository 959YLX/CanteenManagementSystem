package utils

import (
	"github.com/satori/go.uuid"
)

// GenerateToken 生成随机token
func GenerateToken() *string {
	u := uuid.NewV4()
	stringUUID := u.String()
	return &stringUUID
}
