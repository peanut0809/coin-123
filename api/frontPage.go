package api

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/api"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"meta_launchpad/service"
)

type frontPage struct {
	api.CommonBase
}

var FrontPage = new(frontPage)

func (c *frontPage) TransactionSlip(r *ghttp.Request) {
	publisherId := c.GetPublisherId(r)
	transactionSlip, sum := service.FrontPage.TransactionSlip(publisherId)
	req := g.Map{
		"list": transactionSlip,
		"sum":  sum,
	}
	c.SusJsonExit(r, req)
}

func (c *frontPage) VolumeOfTrade(r *ghttp.Request) {
	publisherId := c.GetPublisherId(r)
	count := service.FrontPage.VolumeOfTrade(publisherId)
	req := g.Map{
		"list": count,
	}
	c.SusJsonExit(r, req)
}
