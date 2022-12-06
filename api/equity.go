package api

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/api"
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/client"
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"meta_launchpad/model"
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
	activityId := r.GetInt("id", 1)
	if activityId < 0 {
		c.FailJsonExit(r, "活动id错误")
		return
	}
	ret, err := service.Equity.Info(activityId)
	if err != nil {
		c.FailJsonExit(r, err.Error())
		return
	}
	c.SusJsonExit(r, ret)
}

// CreateOrder 下单
func (c *equity) CreateOrder(r *ghttp.Request) {
	var req model.EquityOrderReq
	err := r.Parse(&req)
	if err != nil {
		c.FailJsonExit(r, err.Error())
		return
	}
	if req.Num <= 0 || req.Id <= 0 {
		c.FailJsonExit(r, "参数错误")
		return
	}
	req.OrderNo = fmt.Sprintf("%d", utils.GetOrderNo())
	req.UserId = c.GetUserId(r)
	req.ClientIp = r.GetClientIp()
	req.PlatformAppId = c.GetAppid(r)
	req.PublisherId = c.GetPublisherId(r)
	queueName := "launchpad.equity"
	mqClient, err := client.GetQueue(client.QueueConfig{
		QueueName:  queueName,
		Exchange:   queueName,
		RoutingKey: "",
		Kind:       "fanout",
		MqUrl:      g.Cfg().GetString("rabbitmq.default.link"),
	})
	if err != nil {
		c.FailJsonExit(r, err.Error())
		return
	}
	defer mqClient.Close()
	err = mqClient.Publish(req)
	if err != nil {
		c.FailJsonExit(r, err.Error())
		return
	}
	c.SusJsonExit(r, req.OrderNo)
}
