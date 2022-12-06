package service

import (
	"github.com/gogf/gf/database/gdb"
	"meta_launchpad/model"
)

type equityOrder struct {
}

var EquityOrder = new(equityOrder)

// 创建订单
func (c *equityOrder) Create(tx *gdb.TX, req model.EquityOrder) (err error) {
	_, err = tx.Model("equity_orders").Insert(&req)
	if err != nil {
		return
	}
	return
}
