package service

import (
	"github.com/gogf/gf/database/gdb"
	"meta_launchpad/model"
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
