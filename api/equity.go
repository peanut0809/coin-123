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
	"time"
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
	// 活动详情
	activityInfo, err := service.Equity.GetValidDetail(req.Id)
	currentTime := time.Now().Unix()
	if currentTime < activityInfo.ActivityStartTime.Unix() {
		c.FailJsonExit(r, "暂未开始")
		return
	}
	if currentTime > activityInfo.ActivityEndTime.Unix() {
		c.FailJsonExit(r, "已结束")
		return
	}
	// 库存判断
	if activityInfo.Number < req.Num {
		c.FailJsonExit(r, "库存不足")
		return
	}
	// 定义限购数量
	limitNum := 0
	// 判断白名单
	if activityInfo.LimitType == model.EQUITY_ACTIVITY_LIMIT_TYPE2 {
		var user *model.EquityUser
		err := g.DB().Model("equity_user").
			Where("activity_id = ?", req.Id).
			Where("user_id = ?", req.UserId).
			Scan(&user)
		if err != nil {
			c.FailJsonExit(r, "网络繁忙")
			return
		}
		if user == nil {
			c.FailJsonExit(r, "不在限购白名单中")
			return
		}
		limitNum = user.LimitNum
	} else {
		limitNum = activityInfo.LimitBuy
	}
	// 判断购买数量
	alreadyBuyNum, err := g.DB().Model("equity_orders").
		Where("user_id = ?", req.UserId).
		Where("activity_id = ?", req.Id).
		Count()
	if err != nil {
		c.FailJsonExit(r, "网络繁忙")
		return
	}
	if alreadyBuyNum >= limitNum {
		c.FailJsonExit(r, "超过限定购买数量")
		return
	}
	// 发送消息
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

// CancelOrder 取消订单
func (c *equity) CancelOrder(r *ghttp.Request) {
	userId := c.GetUserId(r)
	orderNo := r.GetString("orderNo")
	if orderNo == "" {
		c.FailJsonExit(r, "参数错误")
		return
	}
	err := service.EquityOrder.Cancel(userId, orderNo)
	if err != nil {
		c.FailJsonExit(r, err.Error())
		return
	}
	c.SusJsonExit(r)
}

// GetCreateOrderResult 获取下单结果
func (c *equity) GetCreateOrderResult(r *ghttp.Request) {
	orderNo := r.GetQueryString("orderNo")
	ret, err := service.EquityOrder.GetSubResult(orderNo)
	if err != nil {
		return
	}
	c.SusJsonExit(r, ret)
}

// GetOrderList 获取订单列表
func (c *equity) GetOrderList(r *ghttp.Request) {
	pageNum := r.GetQueryInt("pageNum")
	if pageNum <= 0 {
		pageNum = 1
	}
	status := r.GetQueryInt("status")
	userId := c.GetUserId(r)
	publisherId := c.GetPublisherId(r)
	ret, err := service.EquityOrder.GetOrderList(pageNum, userId, status, "", publisherId)
	if err != nil {
		c.FailJsonExit(r, err.Error())
		return
	}
	c.SusJsonExit(r, ret.List)
}

// GetOrderDetail 获取订单详情
func (c *equity) GetOrderDetail(r *ghttp.Request) {
	orderNo := r.GetQueryString("orderNo")
	userId := c.GetUserId(r)
	publisherId := c.GetPublisherId(r)
	ret, err := service.EquityOrder.GetOrderList(1, userId, 0, orderNo, publisherId)
	if err != nil {
		c.FailJsonExit(r, err.Error())
		return
	}
	if len(ret.List) == 0 {
		c.FailJsonExit(r, "订单不存在")
		return
	}
	c.SusJsonExit(r, ret.List[0])
}
