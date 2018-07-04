package service

import (
	"errors"
	"time"

	"geekylx.com/CanteenManagementSystemBackend/src/pb"
	"github.com/golang/protobuf/proto"

	"geekylx.com/CanteenManagementSystemBackend/src/cache"

	"geekylx.com/CanteenManagementSystemBackend/src/database"
	"geekylx.com/CanteenManagementSystemBackend/src/utils"
)

type UserRole uint32
type UserType uint8

const (
	USER_TYPE_ROOT               UserType = 1
	USER_TYPE_ADMIN              UserType = 2
	USER_TYPE_NORMAL             UserType = 3
	USER_TYPE_CANTEEN            UserType = 4
	TOKEN_TTL                    int64    = int64(10 * 60)
	ROLE_CREATE_NORMAL_USER      UserRole = 1
	ROLE_DELETE_NORMAL_USER      UserRole = 1 << 1
	ROLE_CREATE_ADMIN_USER       UserRole = 1 << 2
	ROLE_DELETE_ADMIN_USER       UserRole = 1 << 3
	ROLE_CREATE_CANTEEN          UserRole = 1 << 4
	ROLE_DELETE_CANTEEN          UserRole = 1 << 5
	ROLE_ADD_GOODS               UserRole = 1 << 6
	ROLE_REMOVE_GOODS            UserRole = 1 << 7
	ROLE_RECHARGE                UserRole = 1 << 8
	ROLE_ACCEPT_RECHARGE         UserRole = 1 << 9
	ROLE_CONSUME                 UserRole = 1 << 10
	ROLE_ACCEPT_CONSUME          UserRole = 1 << 11
	ROLE_TRANSFER_ACCOUNT        UserRole = 1 << 12
	ROLE_ACCEPT_TRANSFER_ACCOUNT UserRole = 1 << 13
	ROLE_GET_GOODS_LIST          UserRole = 1 << 14
)

var (
	ErrorPassword   = errors.New("account or password error")
	IllegalArgument = errors.New("argument is illegal")
	ErrorToken      = errors.New("not login or token is outtime")
	ErrorRole       = errors.New("permission denied")
	SystemError     = errors.New("system error")
	TypeDefaultRole = map[UserType]UserRole{
		USER_TYPE_ROOT:    UserRole(^(ROLE_RECHARGE | ROLE_ACCEPT_RECHARGE | ROLE_CONSUME | ROLE_ACCEPT_CONSUME | ROLE_TRANSFER_ACCOUNT | ROLE_ACCEPT_TRANSFER_ACCOUNT)),
		USER_TYPE_ADMIN:   (ROLE_CREATE_NORMAL_USER | ROLE_DELETE_NORMAL_USER | ROLE_RECHARGE),
		USER_TYPE_NORMAL:  UserRole(ROLE_ACCEPT_RECHARGE | ROLE_CONSUME | ROLE_TRANSFER_ACCOUNT | ROLE_ACCEPT_TRANSFER_ACCOUNT | ROLE_GET_GOODS_LIST),
		USER_TYPE_CANTEEN: (ROLE_ADD_GOODS | ROLE_REMOVE_GOODS | ROLE_ACCEPT_CONSUME | ROLE_GET_GOODS_LIST),
	}
)

type userRecordDetail struct {
	Time    time.Time `json:"time"`
	Type    uint8     `json:"type"`
	Account string    `json:"account"`
	Species uint64    `json:"species"`
	GoodsID uint64    `json:"goods"`
	Money   float64   `json:"money"`
}

type canteenRecordDetail struct {
	Time    time.Time `json:"time"`
	Account string    `json:"account"`
	Species uint64    `json:"species"`
	Money   float64   `json:"money"`
	GoodsID uint64    `json:"goods"`
}

// Login 登录业务逻辑
func Login(account string, password string) (token *string, accountType uint8, err error) {
	if utils.IsStringEmpty(account) || utils.IsStringEmpty(password) {
		return nil, 0, IllegalArgument
	}
	encodedPassword, err := utils.EncodePassword(password)
	if err != nil || encodedPassword == nil {
		return nil, 0, err
	}
	userLogin, err := database.GetUserLoginByAccount(account)
	if err != nil || userLogin == nil {
		return nil, 0, ErrorPassword
	}
	storagePassword := userLogin.Password
	if storagePassword != *encodedPassword {
		return nil, 0, ErrorPassword
	}
	token = utils.GenerateToken()
	userInfo, err := database.GetUserInfoByAccount(account)
	if userInfo == nil || err != nil {
		return nil, 0, SystemError
	}
	if cache.TokenCache(*token, account, TOKEN_TTL) != nil {
		return nil, 0, err
	}
	accountType = userInfo.Type
	return
}

// CreateUser 创建用户
func CreateUser(token string, password string, accountType uint8, name string) (account *string, err error) {
	operatorAccount, err := cache.GetAndRefreshToken(token, TOKEN_TTL)
	if operatorAccount == nil || err != nil {
		return nil, ErrorToken
	}
	userInfo, err := database.GetUserInfoByAccount(*operatorAccount)
	if userInfo == nil || err != nil {
		return nil, err
	}
	hasPermission := false
	switch UserType(accountType) {
	case USER_TYPE_ADMIN:
		hasPermission = ((UserRole(userInfo.Role) & ROLE_CREATE_ADMIN_USER) == ROLE_CREATE_ADMIN_USER)
	case USER_TYPE_NORMAL:
		hasPermission = ((UserRole(userInfo.Role) & ROLE_CREATE_NORMAL_USER) == ROLE_CREATE_NORMAL_USER)
	case USER_TYPE_CANTEEN:
		hasPermission = ((UserRole(userInfo.Role) & ROLE_CREATE_CANTEEN) == ROLE_CREATE_CANTEEN)
	default:
		return nil, IllegalArgument
	}
	if !hasPermission {
		return nil, ErrorRole
	}
	encodedPassword, err := utils.EncodePassword(password)
	if encodedPassword == nil || err != nil {
		return nil, err
	}
	account = utils.GenerateAccount()
	if account == nil {
		return nil, SystemError
	}
	newUserInfo := database.UserInfo{
		Account:   *account,
		Type:      accountType,
		Role:      uint32(TypeDefaultRole[UserType(accountType)]),
		Remaining: 0.0,
	}
	newUserLogin := database.UserLogin{
		Account:  *account,
		Password: *encodedPassword,
	}
	if extra, err := proto.Marshal(&pb.UserInfoExtra{
		Name: name,
	}); err == nil {
		newUserInfo.Extra = extra
	}
	database.CreateUserInfo(newUserInfo)
	database.CreateUserLogin(newUserLogin)
	return
}

// DeleteUsers 删除用户
func DeleteUsers(token string, accounts []string) (deletedAccount map[string]bool, err error) {
	if utils.IsStringEmpty(token) || len(accounts) == 0 {
		return nil, IllegalArgument
	}
	operatorAccount, err := cache.GetAndRefreshToken(token, TOKEN_TTL)
	if operatorAccount == nil || err != nil {
		return nil, ErrorToken
	}
	operatorUserInfo, err := database.GetUserInfoByAccount(*operatorAccount)
	if operatorUserInfo == nil || err != nil {
		return nil, err
	}
	userInfos, err := database.GetUserInfosByAccounts(accounts)
	if userInfos == nil || err != nil {
		return nil, err
	}
	deletedAccount = make(map[string]bool)
	for _, account := range accounts {
		deletedAccount[account] = false
	}
	var deleteUsersAccount []string
	for account, userInfo := range userInfos {
		switch UserType(userInfo.Type) {
		case USER_TYPE_NORMAL:
			if UserRole(operatorUserInfo.Role)&ROLE_DELETE_NORMAL_USER == 0 {
				continue
			}
		case USER_TYPE_ADMIN:
			if UserRole(operatorUserInfo.Role)&ROLE_DELETE_NORMAL_USER == 0 {
				continue
			}
		case USER_TYPE_CANTEEN:
			if UserRole(operatorUserInfo.Role)&ROLE_DELETE_CANTEEN == 0 {
				continue
			}
		case USER_TYPE_ROOT:
			continue
		}
		deletedAccount[account] = true
		deleteUsersAccount = append(deleteUsersAccount, account)
	}
	tx := database.Transaction()
	if err = database.DeleteUserInfosByAccountsInTransaction(tx, accounts); err != nil {
		return nil, err
	}
	if err = database.DeleteUserLoginsByAccountsInTransaction(tx, accounts); err != nil {
		return nil, err
	}
	tx.Commit()
	cache.RemoveTokens(accounts)
	return
}

// SelectRecord 查询交易记录
func SelectRecord(token string, startTime int64, endTime int64, species uint8) (totalIncome float64, totalPay float64, details interface{}, err error) {
	if utils.IsStringEmpty(token) {
		return 0, 0, nil, IllegalArgument
	}
	operatorAccount, err := cache.GetAndRefreshToken(token, TOKEN_TTL)
	if operatorAccount == nil || err != nil {
		return 0, 0, nil, ErrorToken
	}
	operatorUserInfo, err := database.GetUserInfoByAccount(*operatorAccount)
	if operatorUserInfo == nil || err != nil {
		return 0, 0, nil, err
	}
	if operatorUserInfo.Type == uint8(USER_TYPE_NORMAL) {
		return selectUserRecord(*operatorAccount)
	} else if operatorUserInfo.Type == uint8(USER_TYPE_CANTEEN) {
		totalIncome, details, err := selectCanteenRecord(*operatorAccount)
		totalPay = 0
		return totalIncome, totalPay, details, err
	}
	return 0, 0, nil, ErrorRole
}

func selectUserRecord(account string) (totalIncome float64, totalPay float64, details []*userRecordDetail, err error) {
	records, err := Select(account, 0, 0, []TransactionType{TRANSACTION_TYPE_CONSUME, TRANSACTION_TYPE_RECHARGE, TRANSACTION_TYPE_TRANSFER_ACCOUNT})
	if records == nil || err != nil {
		return 0, 0, nil, err
	}
	for _, record := range records {
		detail := &userRecordDetail{
			Time:    record.CreatedAt,
			Type:    record.Type,
			Species: record.Species,
		}
		switch TransactionType(detail.Type) {
		case TRANSACTION_TYPE_CONSUME:
			detail.Account = record.To
			if record.Extra != nil {
				extra := pb.FlowingWaterExtra{}
				if err := proto.Unmarshal(record.Extra, &extra); err == nil {
					detail.GoodsID = extra.GoodsID
				}
			}
			detail.Money = -record.Money
			totalPay += record.Money
		case TRANSACTION_TYPE_RECHARGE:
			detail.Account = record.From
			detail.Money = record.Money
			totalIncome += record.Money
		case TRANSACTION_TYPE_TRANSFER_ACCOUNT:
			if record.From == account {
				detail.Money = -record.Money
				detail.Account = record.To
				totalPay += record.Money
			} else {
				detail.Money = record.Money
				detail.Account = record.From
				totalIncome += record.Money
			}
		}
		details = append(details, detail)
	}
	return
}

func selectCanteenRecord(account string) (totalIncome float64, details []*canteenRecordDetail, err error) {
	records, err := Select(account, 0, 0, []TransactionType{TRANSACTION_TYPE_CONSUME})
	if records == nil || err != nil {
		return 0, nil, err
	}
	for _, record := range records {
		detail := &canteenRecordDetail{
			Time:    record.CreatedAt,
			Species: record.Species,
			Money:   record.Money,
			Account: record.From,
		}
		if record.Extra != nil {
			extra := pb.FlowingWaterExtra{}
			if err := proto.Unmarshal(record.Extra, &extra); err == nil {
				detail.GoodsID = extra.GoodsID
			}
		}
		details = append(details, detail)
		totalIncome += record.Money
	}
	return
}
