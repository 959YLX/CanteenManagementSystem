package utils

import (
	"fmt"
	"strconv"

	"geekylx.com/CanteenManagementSystemBackend/src/database"
)

// GenerateAccount 生成账号
func GenerateAccount() *string {
	maxStringAccount, err := database.GetMaxAccount()
	if maxStringAccount == nil || err != nil {
		return nil
	}
	maxAccount, err := strconv.ParseUint(*maxStringAccount, 10, 64)
	if err != nil {
		return nil
	}
	maxAccount++
	account := fmt.Sprintf("%0.10d", maxAccount)
	return &account
}
