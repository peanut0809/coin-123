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

type subscribeActivity struct {
	api.CommonBase
}

var SubscribeActivity = new(subscribeActivity)

//func (s *subscribeActivity) GetSubscribeActivityList(r *ghttp.Request) {
//	ret, err := service.SubscribeActivity.GetList()
//	if err != nil {
//		s.FailJsonExit(r, err.Error())
//		return
//	}
//	s.SusJsonExit(r, ret)
//}

func (s *subscribeActivity) GetSubscribeActivityDetail(r *ghttp.Request) {
	publisherId := s.GetPublisherId(r)
	if publisherId == "" {
		publisherId = r.GetQueryString("publisherId")
	}
	if publisherId == "" {
		s.FailJsonExit(r, "缺少发行商")
		return
	}
	userId := s.GetUserId(r)
	alias := r.GetQueryString("alias")
	if alias == "" {
		s.FailJsonExit(r, "参数错误")
		return
	}
	ret, err := service.SubscribeActivity.GetDetail(alias, userId, publisherId)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *subscribeActivity) GetActivityAwardRecord(r *ghttp.Request) {
	alias := r.GetQueryString("alias")
	if alias == "" {
		s.FailJsonExit(r, "参数错误")
		return
	}
	ret, err := service.SubscribeRecord.GetSubscribeAwardRecord(alias)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *subscribeActivity) GetPayInfo(r *ghttp.Request) {
	userId := s.GetUserId(r)

	alias := r.GetQueryString("alias")
	if alias == "" {
		s.FailJsonExit(r, "参数错误")
		return
	}
	ret, _, err := service.SubscribeActivity.GetMaxBuyNum(alias, userId)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *subscribeActivity) SubActivity(r *ghttp.Request) {
	var req model.DoSubReq
	err := r.Parse(&req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	if req.Type != model.TICKET_MONTH && req.Type != model.TICKET_MONEY && req.Type != model.TICKET_CRYSTAL {
		s.FailJsonExit(r, "type参数错误")
		return
	}
	if req.SubNum <= 0 || req.Alias == "" {
		s.FailJsonExit(r, "参数错误")
		return
	}
	orderNo := fmt.Sprintf("%d", utils.GetOrderNo())
	req.UserId = s.GetUserId(r)
	req.ClientIp = r.GetClientIp()
	req.OrderNo = orderNo
	req.PublisherId = s.GetPublisherId(r)
	req.PlatformAppId = s.GetAppid(r)
	if req.Type == model.TICKET_MONEY && (req.ExitRedirectUrl == "" || req.SuccessRedirectUrl == "") {
		s.FailJsonExit(r, "参数错误")
		return
	}
	queueName := "launchpad.sub"
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
	defer mqClient.Close()
	err = mqClient.Publish(req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, orderNo)
}

func (s *subscribeActivity) GetSubActivityResult(r *ghttp.Request) {
	orderNo := r.GetQueryString("orderNo")
	ret, err := service.SubscribeActivity.GetSubResult(orderNo)
	if err != nil {
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *subscribeActivity) GetSubActivityAwardResult(r *ghttp.Request) {
	aid := r.GetQueryInt("aid")
	userId := s.GetUserId(r)
	ret, err := service.SubscribeActivity.GetSubAwardResult(aid, userId)
	if err != nil {
		return
	}
	s.SusJsonExit(r, ret)
}
