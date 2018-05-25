package service

import (
	"errors"
	"sync"

	"geekylx.com/CanteenManagementSystemBackend/src/cache"
	"geekylx.com/CanteenManagementSystemBackend/src/database"
	"geekylx.com/CanteenManagementSystemBackend/src/utils"
)

var (
	ErrorTargetAccount = errors.New("target account error")
	InsufficientMoney  = errors.New("insufficient money")
	lock               = sync.Mutex{}
)

type TransactionType uint8

const (
	TRANSACTION_TYPE_RECHARGE         TransactionType = 1
	TRANSACTION_TYPE_PAY              TransactionType = 1 << 1
	TRANSACTION_TYPE_TRANSFER_ACCOUNT TransactionType = 1 << 2
)

// Recharge 充值
func Recharge(token string, account string, money float64) (bool, float64, error) {
	if utils.IsStringEmpty(token) || utils.IsStringEmpty(account) || money <= 0 {
		return false, 0, IllegalArgument
	}
	operatorAccount, err := cache.GetAndRefreshToken(token, TOKEN_TTL)
	if operatorAccount == nil || err != nil {
		return false, 0, ErrorToken
	}
	operatorUserInfo, err := database.GetUserInfoByAccount(*operatorAccount)
	if operatorUserInfo == nil || err != nil {
		return false, operatorUserInfo.Remaining, err
	}
	if UserRole(operatorUserInfo.Role)&ROLE_RECHARGE == 0 {
		return false, operatorUserInfo.Remaining, ErrorRole
	}
	rechargeRecord := database.FlowingWater{
		From:  *operatorAccount,
		To:    account,
		Type:  uint8(TRANSACTION_TYPE_RECHARGE),
		Money: money,
	}
	return transaction(*operatorAccount, account, true, money, rechargeRecord)
}

func transaction(fromAccount string, toAccount string, isRecharge bool, money float64, flowingWaterRecord database.FlowingWater) (bool, float64, error) {
	if utils.IsStringEmpty(fromAccount) || utils.IsStringEmpty(toAccount) || money < 0 {
		return false, 0, IllegalArgument
	}
	tx := database.Transaction()
	lock.Lock()
	defer lock.Unlock()
	userInfos, err := database.GetUserInfosByAccounts([]string{fromAccount, toAccount})
	if userInfos == nil || err != nil {
		return false, 0, err
	}
	if _, exist := userInfos[toAccount]; !exist {
		return false, userInfos[fromAccount].Remaining, ErrorTargetAccount
	}
	fromRemaining := userInfos[fromAccount].Remaining
	if !isRecharge {
		if fromRemaining < money {
			return false, fromRemaining, InsufficientMoney
		}
		if err = database.UpdateRemainingInTransaction(tx, fromAccount, fromRemaining-money); err != nil {
			return false, fromRemaining, err
		}
	}
	toRemaining := userInfos[toAccount].Remaining + money
	if err = database.UpdateRemainingInTransaction(tx, toAccount, toRemaining); err != nil {
		return false, fromRemaining, err
	}
	if err = database.RecordInTransaction(tx, flowingWaterRecord); err != nil {
		return false, fromRemaining, err
	}
	tx.Commit()
	if isRecharge {
		return true, toRemaining, nil
	}
	return true, fromRemaining, nil
}
