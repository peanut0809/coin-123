package service

import (
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/model"
)

type activity struct {
}

var Activity = new(activity)

func (s *activity) GetByIds(ids []int) (ret []model.Activity) {
	_ = g.DB().Model("activity").Where("id in (?)", ids).Scan(&ret)
	return
}
