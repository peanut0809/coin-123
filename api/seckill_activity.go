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

type seckillActivity struct {
	api.CommonBase
}

var SeckillActivity = new(seckillActivity)

func (s *seckillActivity) GetDetail(r *ghttp.Request) {
	alias := r.GetQueryString("alias")
	ret, err := service.SeckillActivity.GetValidDetail(alias)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	userId := s.GetUserId(r)
	bnumInfo, err := service.SeckillUserBnum.GetDetail(userId, ret.Id)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	if bnumInfo != nil {
		ret.LimitBuy = bnumInfo.CanBuy
	}
	s.SusJsonExit(r, ret)
}

func (s *seckillActivity) CreateOrder(r *ghttp.Request) {
	var req model.DoBuyReq
	err := r.Parse(&req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	if req.Num <= 0 || req.Alias == "" {
		s.FailJsonExit(r, "参数错误")
		return
	}
	req.OrderNo = fmt.Sprintf("%d", utils.GetOrderNo())
	req.UserId = s.GetUserId(r)
	req.ClientIp = r.GetClientIp()
	req.PlatformAppId = s.GetAppid(r)
	req.PublisherId = s.GetPublisherId(r)
	queueName := "launchpad.seckill"
	mqClient, err := client.GetQueue(client.QueueConfig{
		QueueName:  queueName,
		Exchange:   queueName,
		RoutingKey: "",
		Kind:       "fanout",
		MqUrl:      g.Cfg().GetString("rabbitmq.default.link"),
	})
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	err = mqClient.Publish(req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, req.OrderNo)
}

func (s *seckillActivity) GetCreateOrderResult(r *ghttp.Request) {
	orderNo := r.GetQueryString("orderNo")
	ret, err := service.SeckillActivity.GetSubResult(orderNo)
	if err != nil {
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *seckillActivity) GetOrderList(r *ghttp.Request) {
	pageNum := r.GetQueryInt("pageNum")
	if pageNum <= 0 {
		pageNum = 1
	}
	status := r.GetQueryInt("status")
	userId := s.GetUserId(r)
	ret, err := service.SeckillOrder.GetOrderList(pageNum, userId, status, "")
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *seckillActivity) GetOrderDetail(r *ghttp.Request) {
	orderNo := r.GetQueryString("order_no")
	userId := s.GetUserId(r)
	ret, err := service.SeckillOrder.GetOrderList(1, userId, 0, orderNo)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	if len(ret.List) == 0 {
		s.FailJsonExit(r, "订单不存在")
		return
	}
	s.SusJsonExit(r, ret.List[0])
}