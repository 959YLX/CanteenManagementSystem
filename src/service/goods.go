package service

import (
	"geekylx.com/CanteenManagementSystemBackend/src/cache"
	"geekylx.com/CanteenManagementSystemBackend/src/database"
	"geekylx.com/CanteenManagementSystemBackend/src/pb"
	"geekylx.com/CanteenManagementSystemBackend/src/utils"
	"github.com/golang/protobuf/proto"
)

// GoodsInfo 商品信息数据结构
type GoodsInfo struct {
	Name    string  `json:"name"`
	Species uint64  `json:"species"`
	Price   float64 `json:"price"`
	Canteen string  `json:"canteen"`
}

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
		return false, err
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

// GetGoodsList 获取商品列表
func GetGoodsList(token string) ([]*GoodsInfo, error) {
	if utils.IsStringEmpty(token) {
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
	if UserRole(operatorUserInfo.Role)&ROLE_GET_GOODS_LIST == 0 {
		return nil, ErrorRole
	}
	var belongTo *string
	canteenNameRefMap := make(map[string]string)
	if UserType(operatorUserInfo.Type) == USER_TYPE_CANTEEN {
		belongTo = operatorAccount
		if operatorUserInfo.Extra != nil {
			canteenExtra := &pb.UserInfoExtra{}
			if proto.Unmarshal(operatorUserInfo.Extra, canteenExtra) == nil {
				canteenNameRefMap[*operatorAccount] = canteenExtra.GetName()
			}
		} else {
			canteenNameRefMap[*operatorAccount] = *operatorAccount
		}
	}
	dbGoodsInfos, err := database.GetGoodsList(belongTo)
	if dbGoodsInfos == nil || err != nil {
		return nil, err
	}
	goodsInfos := make([]*GoodsInfo, len(dbGoodsInfos))
	for _, dbGoodsInfo := range dbGoodsInfos {
		goodsInfo := &GoodsInfo{
			Species: dbGoodsInfo.Species,
			Price:   dbGoodsInfo.Price,
		}
		if dbGoodsInfo.Extra != nil {
			goodsPb := &pb.GoodsInfoExtra{}
			if err = proto.Unmarshal(dbGoodsInfo.Extra, goodsPb); err == nil {
				goodsInfo.Name = goodsPb.GetName()
			}
		}
		canteenAccount := dbGoodsInfo.BelongTo
		if _, exist := canteenNameRefMap[canteenAccount]; !exist {
			canteenNameRefMap[canteenAccount] = canteenAccount
			canteenInfo, err := database.GetUserInfoByAccount(canteenAccount)
			if canteenInfo != nil && canteenInfo.Extra != nil {
				canteenNamePb := &pb.UserInfoExtra{}
				if err = proto.Unmarshal(canteenInfo.Extra, canteenNamePb); err == nil {
					canteenNameRefMap[canteenAccount] = canteenNamePb.GetName()
				}
			}
		}
		goodsInfo.Canteen = canteenNameRefMap[canteenAccount]
		goodsInfos = append(goodsInfos)
	}
	return goodsInfos, nil
}
