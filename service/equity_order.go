package service

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	userModel "brq5j1d.gfanx.pro/meta_cloud/meta_service/app/user/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"meta_launchpad/cache"
	"meta_launchpad/model"
	"meta_launchpad/provider"
	"time"
)

type equityOrder struct{}

var EquityOrder = new(equityOrder)

// 创建订单SERVER
func (c *equityOrder) Create(req *model.EquityOrderReq, activityInfo *model.EquityActivity) (err error) {
	// 扣除库存
	var tx *gdb.TX
	tx, e := g.DB().Begin()
	if e != nil {
		c.SetSubResult(model.EquitySubResult{
			Reason:  e.Error(),
			Step:    "fail",
			OrderNo: req.OrderNo,
		})
		return
	}
	r, e := tx.Exec("UPDATE equity_activity SET number = number - ? WHERE id = ?", req.Num, activityInfo.Id)
	if e != nil {
		err = tx.Rollback()
		c.SetSubResult(model.EquitySubResult{
			Reason:  e.Error(),
			Step:    "fail",
			OrderNo: req.OrderNo,
		})
		return
	}
	affectedNum, _ := r.RowsAffected()
	if affectedNum != 1 {
		err = tx.Rollback()
		EquityOrder.SetSubResult(model.EquitySubResult{
			Reason:  "更新库存失败",
			Step:    "fail",
			OrderNo: req.OrderNo,
		})
		return
	}
	// 用户信息
	user := &userModel.Users{}
	params := &map[string]interface{}{
		"userId": req.UserId,
	}
	// 获取用户信息
	result, err := utils.SendJsonRpc(context.Background(), "ucenter", "UserBase.GetOneUserInfo", params)
	if err != nil {
		g.Log().Error(err)
		return
	}
	err = json.Unmarshal([]byte(gconv.String(result)), user)
	if err != nil {
		g.Log().Error(err)
		return
	}
	// 生成订单
	_, err = tx.Model("equity_orders").Insert(&model.EquityOrder{
		PublisherId:  req.PublisherId,
		OrderNo:      req.OrderNo,
		Num:          req.Num,
		RealFee:      req.Num * activityInfo.Price,
		ActivityId:   activityInfo.Id,
		ActivityName: activityInfo.Name,
		UserName:     user.Nickname,
		UserId:       req.UserId,
		Status:       model.WAIT_PAY,
		Price:        activityInfo.Price,
		PayExpireAt:  gtime.Now().Add(time.Minute * 10),
	})
	if err != nil {
		err = tx.Rollback()
		return
	}
	// 生成预支付订单
	orderReq := new(provider.CreateOrderReq)
	orderReq.ClientIp = req.ClientIp
	orderReq.UserId = req.UserId
	orderReq.AppType = "launchpad_equity"
	orderReq.PayAmount = activityInfo.Price * req.Num
	orderReq.PayExpire = gtime.Now().Add(time.Minute * 10)
	orderReq.Subject = "权益活动"
	orderReq.Description = "权益活动"
	orderReq.SuccessRedirectUrl = req.SuccessRedirectUrl
	orderReq.ExitRedirectUrl = req.ExitRedirectUrl
	orderReq.AppOrderNo = req.OrderNo
	orderReq.PublisherId = req.PublisherId
	orderReq.PlatformAppId = req.PlatformAppId
	e = provider.Payment.CreateOrder(orderReq)
	if e != nil {
		err = tx.Rollback()
		c.SetSubResult(model.EquitySubResult{
			Reason:  "下单失败",
			Step:    "fail",
			OrderNo: req.OrderNo,
		})
		g.Log().Errorf("equity err:%v", e)
		return
	}
	err = tx.Commit()
	return
}

// 订单列表SERVER
func (c *equityOrder) GetOrderList(pageNum int, userId string, status int, orderNo, publisherId string) (ret model.EquityOrderList, err error) {
	m := g.DB().Model("equity_orders").Where("publisher_id = ? AND user_id = ?", publisherId, userId)
	if status != 0 {
		m = m.Where("status", status)
	}
	if orderNo != "" {
		m = m.Where("order_no", orderNo)
	}
	ret.Total, err = m.Count()
	if err != nil {
		return
	}
	if ret.Total == 0 {
		return
	}
	var list []*model.EquityOrder
	err = m.Order("id DESC").Page(pageNum, 20).Scan(&list)
	if err != nil {
		return
	}
	for _, v := range list {
		lastSec := v.PayExpireAt.Unix() - time.Now().Unix()
		if lastSec <= 0 {
			lastSec = 0
		}
		ret.List = append(ret.List, &model.EquityOrderFull{
			EquityOrder: v,
			PriceYuan:   fmt.Sprintf("%.2f", float64(v.Price)/100),
			RealFeeYuan: fmt.Sprintf("%.2f", float64(v.RealFee)/100),
			LastSec:     lastSec,
		})
	}
	return
}

// RedisSet订单信息SERVER
func (c *equityOrder) SetSubResult(req model.EquitySubResult) {
	_, err := g.Redis().Do("SET", fmt.Sprintf(model.SubSetEquityResultKey, req.OrderNo), gconv.String(req), "ex", 3600)
	if err != nil {
		g.Log().Errorf("EquityBuy err:%v", err)
		return
	}
	return
}

// RedisGet订单信息SERVER
func (c *equityOrder) GetSubResult(orderNo string) (ret model.EquitySubResult, err error) {
	gv, e := g.Redis().DoVar("GET", fmt.Sprintf(model.SubSetEquityResultKey, orderNo))
	if e != nil {
		err = e
		g.Log().Errorf("EquityBuy err:%v", err)
		return
	}
	if gv == nil {
		err = fmt.Errorf("内部错误，请重新下单")
		return
	}
	err = gv.Scan(&ret)
	return
}

// 取消订单SERVER
func (c *equityOrder) Cancel(userId string, orderNo string) (err error) {
	orderInfo, e := c.GetInfoByOrderNo(orderNo)
	if e != nil {
		err = e
		return
	}
	if orderInfo.UserId != userId {
		err = fmt.Errorf("无权操作")
		return
	}
	if orderInfo.Status != 1 {
		err = fmt.Errorf("当前状态不能取消")
		return
	}
	now := time.Now()
	if now.Unix() >= orderInfo.PayExpireAt.Unix() {
		err = fmt.Errorf("订单已过期")
		return
	}
	if orderInfo.PayExpireAt.Unix()-now.Unix() < 300 { //超过5分钟了,算超时
		//设置处罚时间
		_, err := g.Redis().Do("SET", fmt.Sprintf(cache.EQUITY_DISCIPLINE, userId), 1, "ex", 3600*24*30)
		if err != nil {
			return err
		}
		_ = c.UpdateOrderNoStatus(orderNo, model.TIMEOUT)
	}
	err = g.DB().Transaction(context.Background(), func(ctx context.Context, tx *gdb.TX) error {
		err = c.UpdateOrderNoStatus(orderNo, model.CANCEL)
		if err != nil {
			return err
		}
		return c.InventoryRollback(tx, orderInfo.ActivityId, orderInfo.Num)
	})
	return
}

// 更改订单状态
func (c *equityOrder) UpdateOrderNoStatus(orderNo string, status int) (err error) {
	_, err = g.DB().Exec("UPDATE equity_orders SET status = ? WHERE order_no = ?", status, orderNo)
	return
}

// 根据订单号获取订单信息
func (c *equityOrder) GetInfoByOrderNo(orderNo string) (ret model.EquityOrder, err error) {
	err = g.DB().Model("equity_orders").Where("order_no = ?", orderNo).Scan(&ret)
	if err != nil {
		return
	}
	return
}

// 回退库存
func (c *equityOrder) InventoryRollback(tx *gdb.TX, activityId int, num int) (err error) {
	_, err = tx.Exec("UPDATE equity_activity SET number = number + ? WHERE id = ?", num, activityId)
	return
}
