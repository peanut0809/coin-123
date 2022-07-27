package service

import (
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/model"
	"time"
)

var SeckillOrder = new(seckillOrder)

type seckillOrder struct {
}

func (s *seckillOrder) Create(tx *gdb.TX, in model.SeckillOrder) (err error) {
	_, err = tx.Model("seckill_orders").Insert(&in)
	if err != nil {
		return
	}
	return
}

func (s *seckillOrder) GetOrderList(pageNum int, userId string, status int, orderNo string) (ret model.SeckillOrderList, err error) {
	m := g.DB().Model("seckill_orders").Where("user_id", userId)
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
	list := make([]model.SeckillOrder, 0)
	err = m.Page(pageNum, 20).Scan(&list)
	if err != nil {
		return
	}
	for _, v := range list {
		lastSec := v.PayExpireAt.Unix() - time.Now().Unix()
		if lastSec <= 0 {
			lastSec = 0
		}
		ret.List = append(ret.List, model.SeckillOrderFull{
			SeckillOrder: v,
			PriceYuan:    fmt.Sprintf("%.2f", float64(v.Price)/100),
			RealFeeYuan:  fmt.Sprintf("%.2f", float64(v.RealFee)/100),
			LastSec:      lastSec,
		})
	}
	return
}

func (s *seckillOrder) GetByOrderNos(orderNos []string) (ret []model.SeckillOrder, err error) {
	err = g.DB().Model("seckill_orders").Where("order_no in (?)", orderNos).Scan(&ret)
	if err != nil {
		return
	}
	return
}

func (s *seckillOrder) UpdateOrderNosStatus(orderNos []string, status int) (err error) {
	_, err = g.DB().Exec("UPDATE seckill_orders SET status = ? WHERE order_no in (?)", status, orderNos)
	return
}