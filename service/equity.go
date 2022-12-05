package service

import (
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/model"
)

type equity struct{}

var Equity = new(equity)

func (c *equity) List(publisherId string, pageNum int, pageSize int) (res *model.EquityActivityList, err error) {
	m := g.DB().Model("activity")

	if publisherId != "" {
		m = m.Where("publisher_id = ?", publisherId)
	}
	var equity []model.EquityActivity
	err = m.Order("id DESC").Page(pageNum, pageSize).Scan(&equity)
	if err != nil {
		return
	}
	res.Total, err = m.Count()
	if err != nil {
		return
	}
	if res.Total == 0 {
		return
	}
	for _, v := range equity {
		res.List = append(res.List, &v)
	}
	return
}
