package service

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/model"
	"time"
)

var SeckillWaitPayOrder = new(seckillWaitPayOrder)

type seckillWaitPayOrder struct {
}

func (s *seckillWaitPayOrder) Create(tx *gdb.TX, in model.SeckillWaitPayOrder) (err error) {
	_, err = tx.Model("seckill_wait_pay_orders").Insert(&in)
	if err != nil {
		return
	}
	return
}

func (s *seckillWaitPayOrder) GetWaitPayOrder() (ret []model.SeckillWaitPayOrder, err error) {
	err = g.DB().Model("seckill_wait_pay_orders").Where("? > pay_expire_at", time.Now()).Scan(&ret)
	if err != nil {
		return
	}
	return
}
