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
	ErrorPassword   = errors.New("password error")
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
		return nil, 0, err
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
		return nil, 0, err
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
