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

// VolumeOfTrade 近期支付数
func (c *frontPage) VolumeOfTrade(r *ghttp.Request) {
	publisherId := c.GetPublisherId(r)
	day := r.GetInt("day")
	dealTime, paymentTime, dealCount, paymentCount := service.FrontPage.VolumeOfTrade(publisherId, day)
	req := g.Map{
		"dealNum": g.Map{
			"dealTime":  dealTime,
			"dealCount": dealCount,
		},
		"payment": g.Map{
			"paymentTime":  paymentTime,
			"paymentCount": paymentCount,
		},
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
	priceFloat, priceTime, float, count, err := service.FrontPage.Turnover(publisherId)
	if err != nil {
		return
	}
	req := g.Map{
		"list": g.Map{
			"priceFloat": priceFloat,
			"priceTime":  priceTime,
		},
		"float": float,
		"count": count,
	}
	c.SusJsonExit(r, req)
}

// 免费版

func (c *frontPage) FreeTransactionSlip(r *ghttp.Request) {
	userId := c.GetUserId(r)
	transactionSlip, sum := service.FrontPage.FreeTransactionSlip(userId)
	req := g.Map{
		"list": transactionSlip,
		"sum":  sum,
	}
	c.SusJsonExit(r, req)
}

// VolumeOfTrade 近期支付数
func (c *frontPage) FreeVolumeOfTrade(r *ghttp.Request) {
	userId := c.GetUserId(r)
	day := r.GetInt("day")
	dealTime, paymentTime, dealCount, paymentCount := service.FrontPage.FreeVolumeOfTrade(userId, day)
	req := g.Map{
		"dealNum": g.Map{
			"dealTime":  dealTime,
			"dealCount": dealCount,
		},
		"payment": g.Map{
			"paymentTime":  paymentTime,
			"paymentCount": paymentCount,
		},
	}
	c.SusJsonExit(r, req)
}

func (c *frontPage) FreeTransactionsNum(r *ghttp.Request) {
	userId := c.GetUserId(r)
	count, float, err := service.FrontPage.FreeTransactionsNum(userId)
	if err != nil {
		return
	}
	req := g.Map{
		"count": count,
		"float": float,
	}
	c.SusJsonExit(r, req)
}

func (c *frontPage) FreePayers(r *ghttp.Request) {
	userId := c.GetUserId(r)
	count, float, err := service.FrontPage.FreePayers(userId)
	if err != nil {
		return
	}
	req := g.Map{
		"count": count,
		"float": float,
	}
	c.SusJsonExit(r, req)
}

func (c *frontPage) FreeTurnover(r *ghttp.Request) {
	userId := c.GetUserId(r)
	priceFloat, priceTime, float, count, err := service.FrontPage.FreeTurnover(userId)
	if err != nil {
		return
	}
	req := g.Map{
		"list": g.Map{
			"priceFloat": priceFloat,
			"priceTime":  priceTime,
		},
		"float": float,
		"count": count,
	}
	c.SusJsonExit(r, req)
}
