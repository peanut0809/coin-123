package service

import (
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"meta_launchpad/model"
	"time"
)

const WAIT_PAY = 1 // 待支付

type equityOrder struct{}

var EquityOrder = new(equityOrder)

// 创建订单
func (c *equityOrder) Create(tx *gdb.TX, req model.EquityOrder) (err error) {
	_, err = tx.Model("equity_orders").Insert(&req)
	if err != nil {
		return
	}
	return
}

// 订单列表
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

// RedisSet订单信息
func (c *equityOrder) SetSubResult(req model.EquitySubResult) {
	_, err := g.Redis().Do("SET", fmt.Sprintf(model.SubSetEquityResultKey, req.OrderNo), gconv.String(req), "ex", 3600)
	if err != nil {
		g.Log().Errorf("EquityBuy err:%v", err)
		return
	}
	return
}

// RedisGet订单信息
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
