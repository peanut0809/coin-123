package api

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/api"
	"github.com/gogf/gf/net/ghttp"
	"meta_launchpad/service"
)

type equity struct {
	api.CommonBase
}

var Equity = new(equity)

// List 活动列表
func (c *equity) List(r *ghttp.Request) {
	pageNum := r.GetInt("pageNum", 1)
	pageSize := r.GetInt("pageSize", 20)
	publisherId := r.GetString("publisherId")
	ret, err := service.Equity.List(publisherId, pageNum, pageSize)
	if err != nil {
		c.FailJsonExit(r, err.Error())
		return
	}
	c.SusJsonExit(r, ret)
}

// Info 活动详情
func (c *equity) Info(r *ghttp.Request) {
	//activityId := r.GetInt("activity_id", 1)
}
