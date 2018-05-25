package service

import (
	"errors"
	"sync"

	"geekylx.com/CanteenManagementSystemBackend/src/pb"
	"github.com/golang/protobuf/proto"

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
	TRANSACTION_TYPE_CONSUME          TransactionType = 1 << 1
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
	userInfos, err := database.GetUserInfosByAccounts([]string{*operatorAccount, account})
	if len(userInfos) != 2 {
		return false, 0, ErrorTargetAccount
	}
	if UserRole(userInfos[*operatorAccount].Role)&ROLE_RECHARGE == 0 || UserRole(userInfos[account].Role)&ROLE_ACCEPT_RECHARGE == 0 {
		return false, userInfos[account].Remaining, ErrorRole
	}
	rechargeRecord := database.FlowingWater{
		From:  *operatorAccount,
		To:    account,
		Type:  uint8(TRANSACTION_TYPE_RECHARGE),
		Money: money,
	}
	return transaction(*operatorAccount, account, true, money, rechargeRecord)
}

// Consume 消费接口
func Consume(token string, fromAccount string, goodsID uint) (bool, float64, error) {
	if utils.IsStringEmpty(token) || utils.IsStringEmpty(fromAccount) {
		return false, 0, IllegalArgument
	}
	operatorAccount, err := cache.GetAndRefreshToken(token, TOKEN_TTL)
	if operatorAccount == nil || err != nil {
		return false, 0, ErrorToken
	}
	userInfos, err := database.GetUserInfosByAccounts([]string{*operatorAccount, fromAccount})
	if len(userInfos) != 2 {
		return false, 0, ErrorTargetAccount
	}
	fromRemaining := userInfos[fromAccount].Remaining
	if UserRole(userInfos[fromAccount].Role)&ROLE_CONSUME == 0 || UserRole(userInfos[*operatorAccount].Role)&ROLE_ACCEPT_CONSUME == 0 {
		return false, fromRemaining, ErrorRole
	}
	goodsInfo, err := database.GetGoodsInfo(goodsID)
	if goodsInfo == nil || goodsInfo.BelongTo != *operatorAccount || err != nil {
		return false, fromRemaining, IllegalArgument
	}
	extraPb, err := proto.Marshal(&pb.FlowingWaterExtra{
		GoodsID: uint64(goodsID),
	})
	if err != nil {
		return false, fromRemaining, err
	}
	consumeRecord := database.FlowingWater{
		From:    fromAccount,
		To:      *operatorAccount,
		Type:    uint8(TRANSACTION_TYPE_CONSUME),
		Money:   goodsInfo.Price,
		Species: goodsInfo.Species,
		Extra:   extraPb,
	}
	return transaction(fromAccount, *operatorAccount, false, goodsInfo.Price, consumeRecord)
}

// TransferAccount 转账业务
func TransferAccount(token string, toAccount string, money float64) (bool, float64, error) {
	if utils.IsStringEmpty(token) || utils.IsStringEmpty(toAccount) || money < 0 {
		return false, 0, IllegalArgument
	}
	operatorAccount, err := cache.GetAndRefreshToken(token, TOKEN_TTL)
	if operatorAccount == nil || err != nil {
		return false, 0, ErrorToken
	}
	if *operatorAccount == toAccount {
		return false, 0, ErrorTargetAccount
	}
	userInfos, err := database.GetUserInfosByAccounts([]string{*operatorAccount, toAccount})
	if userInfos == nil || err != nil {
		return false, 0, ErrorTargetAccount
	}
	if UserRole(userInfos[*operatorAccount].Role)&ROLE_TRANSFER_ACCOUNT == 0 || UserRole(userInfos[toAccount].Role)&ROLE_ACCEPT_TRANSFER_ACCOUNT == 0 {
		return false, 0, ErrorRole
	}
	transferAccountRecord := database.FlowingWater{
		From:  *operatorAccount,
		To:    toAccount,
		Type:  uint8(TRANSACTION_TYPE_TRANSFER_ACCOUNT),
		Money: money,
	}
	return transaction(*operatorAccount, toAccount, false, money, transferAccountRecord)
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
	return true, fromRemaining - money, nil
}
