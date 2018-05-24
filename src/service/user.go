package service

import (
	"errors"

	"geekylx.com/CanteenManagementSystemBackend/src/cache"

	"geekylx.com/CanteenManagementSystemBackend/src/database"
	"geekylx.com/CanteenManagementSystemBackend/src/utils"
)

type UserRole uint8
type UserType uint8

const (
	USER_TYPE_ROOT          UserType = 1
	USER_TYPE_ADMIN         UserType = 2
	USER_TYPE_NORMAL        UserType = 3
	TOKEN_TTL               int64    = int64(10 * 60)
	ROLE_CREATE_NORMAL_USER UserRole = 1
	ROLE_DELETE_NORMAL_USER UserRole = 1 << 1
	ROLE_CREATE_ADMIN_USER  UserRole = 1 << 2
	ROLE_DELETE_ADMIN_USER  UserRole = 1 << 3
)

var (
	ErrorPassword   = errors.New("account or password error")
	IllegalArgument = errors.New("argument is illegal")
	ErrorToken      = errors.New("not login or token is outtime")
	ErrorRole       = errors.New("permission denied")
	SystemError     = errors.New("system error")
	TypeDefaultRole = map[UserType]UserRole{
		USER_TYPE_ROOT:   UserRole(255),
		USER_TYPE_ADMIN:  (ROLE_CREATE_NORMAL_USER | ROLE_DELETE_NORMAL_USER),
		USER_TYPE_NORMAL: UserRole(0),
	}
)

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
	token, err = utils.GenerateToken()
	if token == nil || err != nil {
		return nil, 0, err
	}
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
func CreateUser(token string, password string, accountType uint8) (account *string, err error) {
	operatorAccount, err := cache.GetAndRefreshToken(token, TOKEN_TTL)
	if operatorAccount == nil || err != nil {
		return nil, ErrorToken
	}
	userLogin, err := database.GetUserLoginByAccount(*operatorAccount)
	if userLogin == nil || err != nil {
		return nil, err
	}
	hasPermission := false
	switch UserType(accountType) {
	case USER_TYPE_ADMIN:
		hasPermission = ((UserRole(userLogin.Role) & ROLE_CREATE_ADMIN_USER) == ROLE_CREATE_ADMIN_USER)
	case USER_TYPE_NORMAL:
		hasPermission = ((UserRole(userLogin.Role) & ROLE_CREATE_NORMAL_USER) == ROLE_CREATE_NORMAL_USER)
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
		Remaining: 0.0,
	}
	newUserLogin := database.UserLogin{
		Account:  *account,
		Role:     uint8(TypeDefaultRole[UserType(accountType)]),
		Password: *encodedPassword,
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
	userLogin, err := database.GetUserLoginByAccount(*operatorAccount)
	if userLogin == nil || err != nil {
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
	for _, userInfo := range userInfos {
		switch UserType(userInfo.Type) {
		case USER_TYPE_NORMAL:
			if UserRole(userLogin.Role)&ROLE_DELETE_NORMAL_USER == 0 {
				continue
			}
		case USER_TYPE_ADMIN:
			if UserRole(userLogin.Role)&ROLE_DELETE_NORMAL_USER == 0 {
				continue
			}
		case USER_TYPE_ROOT:
			continue
		}
		deletedAccount[userInfo.Account] = true
		deleteUsersAccount = append(deleteUsersAccount, userInfo.Account)
	}
	tx := database.Transaction()
	if err = database.DeleteUserInfosByAccountsInTransaction(tx, accounts); err != nil {
		return nil, err
	}
	if err = database.DeleteUserLoginsByAccountsInTransaction(tx, accounts); err != nil {
		return nil, err
	}
	tx.Commit()
	return
}
