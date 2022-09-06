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
	day := r.GetInt("day")
	dealNum, payment := service.FrontPage.VolumeOfTrade(publisherId, day)
	req := g.Map{
		"dealNum": dealNum,
		"payment": payment,
	}
	c.SusJsonExit(r, req)
}

func (c *frontPage) TransactionsNum(r *ghttp.Request) {
	publisherId := c.GetPublisherId(r)
	count, float, err := service.FrontPage.TransactionsNum(publisherId)
	if err != nil {
		return
	}
	req := g.Map{
		"count": count,
		"float": float,
	}
	c.SusJsonExit(r, req)
}

func (c *frontPage) Payers(r *ghttp.Request) {
	publisherId := c.GetPublisherId(r)
	count, float, err := service.FrontPage.Payers(publisherId)
	if err != nil {
		return
	}
	req := g.Map{
		"count": count,
		"float": float,
	}
	c.SusJsonExit(r, req)
}

func (c *frontPage) Turnover(r *ghttp.Request) {
	publisherId := c.GetPublisherId(r)
	price, float, count, err := service.FrontPage.Turnover(publisherId)
	if err != nil {
		return
	}
	req := g.Map{
		"list":  price,
		"float": float,
		"count": count,
	}
	c.SusJsonExit(r, req)
}
