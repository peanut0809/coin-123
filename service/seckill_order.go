package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/cache"
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

func (s *seckillOrder) Count(id int) (count int) {
	count, _ = g.DB().Model("seckill_orders").Where("aid = ?", id).Count()
	return
}

func (s *seckillOrder) GetOrderList(pageNum int, userId string, status int, orderNo, publisherId string) (ret model.SeckillOrderList, err error) {
	m := g.DB().Model("seckill_orders").Where("user_id = ? AND publisher_id = ?", userId, publisherId)
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
	err = m.Order("id DESC").Page(pageNum, 20).Scan(&list)
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

func (s *seckillOrder) Cancel(userId string, orderNo string) (err error) {
	orderInfo, e := s.GetByOrderNos([]string{orderNo})
	if e != nil {
		err = e
		return
	}
	if len(orderInfo) == 0 {
		err = fmt.Errorf("订单不存在")
		return
	}
	if orderInfo[0].UserId != userId {
		err = fmt.Errorf("无权操作")
		return
	}
	if orderInfo[0].Status != 1 {
		err = fmt.Errorf("当前状态不能取消")
		return
	}
	now := time.Now()
	if now.Unix() >= orderInfo[0].PayExpireAt.Unix() {
		err = fmt.Errorf("订单已过期")
		return
	}
	if orderInfo[0].PayExpireAt.Unix()-now.Unix() < 300 { //超过5分钟了,算超时
		//设置处罚时间
		g.Redis().Do("SET", fmt.Sprintf(cache.SECKILL_DISCIPLINE, userId), 1, "ex", 3600*24*30)
		_ = SeckillOrder.UpdateOrderNosStatus([]string{orderNo}, 3)
	}
	_ = SeckillOrder.UpdateOrderNosStatus([]string{orderNo}, 4)
	err = g.DB().Transaction(context.Background(), func(ctx context.Context, tx *gdb.TX) error {
		return s.CancelHandel(tx, orderInfo[0].Aid, orderInfo[0].Num, orderNo, userId)
	})
	return
}

func (s *seckillOrder) CancelHandel(tx *gdb.TX, aid int, num int, orderNo string, userId string) (err error) {
	err = SeckillActivity.UpdateRemain(tx, aid, num)
	if err != nil {
		return
	}
	err = SeckillUserBnum.UpdateRemain(tx, userId, aid, num)
	if err != nil {
		return
	}
	//删除待支付的订单
	err = SeckillWaitPayOrder.Del(orderNo)
	if err != nil {
		return
	}
	return
}
