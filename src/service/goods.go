package service

import (
	"geekylx.com/CanteenManagementSystemBackend/src/cache"
	"geekylx.com/CanteenManagementSystemBackend/src/database"
	"geekylx.com/CanteenManagementSystemBackend/src/pb"
	"geekylx.com/CanteenManagementSystemBackend/src/utils"
	"github.com/golang/protobuf/proto"
)

// AddGoods 新增商品接口
func AddGoods(token string, species uint64, price float64, name string) (bool, error) {
	if utils.IsStringEmpty(token) {
		return false, IllegalArgument
	}
	operatorAccount, err := cache.GetAndRefreshToken(token, TOKEN_TTL)
	if operatorAccount == nil || err != nil {
		return false, ErrorToken
	}
	operatorUserInfo, err := database.GetUserInfoByAccount(*operatorAccount)
	if operatorUserInfo == nil || err != nil {
		return false, SystemError
	}
	if UserType(operatorUserInfo.Type) != USER_TYPE_CANTEEN || UserRole(operatorUserInfo.Role)&ROLE_ADD_GOODS == 0 {
		return false, ErrorRole
	}
	extraPb, err := proto.Marshal(&pb.GoodsInfoExtra{
		Name: name,
	})
	if err != nil {
		return false, err
	}
	goods := database.GoodsInfo{
		BelongTo: *operatorAccount,
		Species:  species,
		Price:    price,
		Extra:    extraPb,
	}
	database.SaveGoods(goods)
	return true, nil
}
