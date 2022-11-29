package provider

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"strings"
)

type user struct {
}

var User = new(user)

type GetUserInfoRet struct {
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	UserId   string `json:"userId"`
	Crystal  int    `json:"crystal"`
}

func (s *user) GetUserInfo(userIds []string) (ret []GetUserInfoRet, retMap map[string]GetUserInfoRet, err error) {
	params := &map[string]interface{}{
		"userIds": userIds,
	}
	result, e := utils.SendJsonRpc(context.Background(), "ucenter", "UserBase.GetUserInfo", params)
	if e != nil {
		err = e
		return
	}
	err = json.Unmarshal([]byte(gconv.String(result)), &ret)
	if err != nil {
		return
	}
	retMap = make(map[string]GetUserInfoRet)
	for _, v := range ret {
		retMap[v.UserId] = v
	}
	return
}

type UserMemberLevel struct {
	Id          int64       `orm:"id,primary" json:"id"`            //
	UserId      string      `orm:"user_id" json:"userId"`           // 用户ID
	LevelCode   string      `orm:"level_code" json:"levelCode"`     // 身份标识
	ExpiredTime *gtime.Time `orm:"expired_time" json:"expiredTime"` // 过期时间
	CreateTime  *gtime.Time `orm:"create_time" json:"createTime"`   // 创建时间
	Source      string      `orm:"source" json:"source"`            // 来源
	Status      int         `orm:"status" json:"status"`            // 状态
	UpdateTime  *gtime.Time `orm:"update_time" json:"updateTime"`   // 修改时间
	Scene       string      `orm:"scene" json:"scene"`
}

//获取用户会员信息
func (s *user) GetUserMemberLevel(userId, scene string) (ret *UserMemberLevel, err error) {
	params := &map[string]interface{}{
		"userId": userId,
		"scene":  scene,
	}
	err = utils.SendJsonRpcScan(context.Background(), "ucenter", "Users.GetUserMemberLevel", params, &ret)
	if err != nil {
		g.Log().Errorf("rpc Users.GetUserMemberLevel err:%v", err)
		return
	}
	return
}

//封号
type UserSealupReq struct {
	Id        int         `json:"id"`
	UserId    string      `json:"userId"`
	IsForever int         `json:"isForever"`
	StartAt   *gtime.Time `json:"startAt"`
	EndAt     *gtime.Time `json:"endAt"`
	Reason    string      `json:"reason"`
	Scene     string      `json:"scene"`
}

func (s *user) SealupUser(in *UserSealupReq) (err error) {
	var ret interface{}
	err = utils.SendJsonRpcScan(context.Background(), "ucenter", "UserBase.SealupUser", in, ret)
	if err != nil {
		g.Log().Errorf("rpc UserBase.SealupUser err:%v", err)
		return
	}
	return
}

func (s *user) Logout(userId string) (err error) {
	var ret interface{}
	params := &map[string]interface{}{
		"userId":    userId,
		"loginFrom": "MARKET",
	}
	err = utils.SendJsonRpcScan(context.Background(), "ucenter", "UserBase.Logout", params, ret)
	if err != nil {
		g.Log().Errorf("rpc UserBase.Logout err:%v", err)
		return
	}
	return
}

type GetUserMonthTicketRet struct {
	MonthTicket int `json:"monthTicket"`
}

//获取用户月票
func (s *user) GetUserMonthTicket(userId string) (ret *GetUserMonthTicketRet, err error) {
	params := &map[string]interface{}{
		"userId": userId, //用户ID
	}
	err = utils.SendJsonRpcScan(context.Background(), "ucenter", "MonthTicket.GetUserMonthTicket", params, &ret)
	if err != nil {
		g.Log().Errorf("rpc Users.GetUserMonthTicket err:%v", err)
		return
	}
	return
}

//操作月票
func (s *user) OptUserMonthTicket(params *map[string]interface{}) (err error) {
	var ret interface{}
	err = utils.SendJsonRpcScan(context.Background(), "ucenter", "MonthTicket.OptUserMonthTicket", params, ret)
	if err != nil {
		if strings.Contains(err.Error(), "monthTicket not enough") {
			err = fmt.Errorf("月票数量不足")
			return
		}
		g.Log().Errorf("OptUserMonthTicket err:%v", err)
		return
	}
	return
}

//操作元晶
func (s *user) YJTransfer(params *map[string]interface{}) (err error) {
	var ret interface{}
	err = utils.SendJsonRpcScan(context.Background(), "ucenter", "Users.YJTransfer", params, ret)
	if err != nil {
		if strings.Contains(err.Error(), "余额不足") {
			err = fmt.Errorf("元晶余额不足")
			return
		}
		g.Log().Errorf("OptUserMonthTicket err:%v", err)
		return
	}
	return
}

func (s *user) GetUserInfoByPhone(params *map[string]interface{}) (ret map[string]GetUserInfoRet, err error) {
	err = utils.SendJsonRpcScan(context.Background(), "ucenter", "Users.GetUserInfoByPhone", params, &ret)
	if err != nil {
		g.Log().Errorf("GetUserInfoByPhone err:%v", err)
		return
	}
	return
}

type GetStoreBalanceRes struct {
	Id          int        `json:"id"`          // ID
	PublisherId string     `json:"publisherId"` // 发行商ID
	Balance     int        `json:"balance"`     // 余额
	CreatedAt   gtime.Time `json:"createdAt"`   // 创建时间
	UpdatedAt   gtime.Time `json:"updatedAt"`   // 更新时间
	TotalAmount int        `json:"totalAmount"` // 累计充值
	UserId      string     `json:"userId"`      // 用户ID
	PayPassword string     `json:"payPassword"` // 支付密码
}

// GetStoreBalance 获取店铺余额
func (s *user) GetStoreBalance(userId string, publisherId string) (res *GetStoreBalanceRes, err error) {
	params := g.Map{
		"userId":      userId,
		"publisherId": publisherId,
	}
	err = utils.SendJsonRpcScan(context.Background(), "ucenter", "CnyPublisher.GetPublisherAccount", params, res)
	if err != nil {
		g.Log().Errorf("GetUserInfoByPhone err:%v", err)
		return
	}
	return
}
