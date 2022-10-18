package api

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/api"
	"github.com/gogf/gf/net/ghttp"
	"meta_launchpad/model"
	"meta_launchpad/service"
)

type subscribeRecord struct {
	api.CommonBase
}

var SubscribeRecord = new(subscribeRecord)

func (s *subscribeRecord) GetList(r *ghttp.Request) {
	publisherId := s.GetPublisherId(r)
	if publisherId == "" {
		publisherId = r.GetQueryString("publisherId")
	}
	userId := s.GetUserId(r)
	pageNum := r.GetQueryInt("pageNum", 1)
	if pageNum <= 0 {
		pageNum = 1
	}
	award := r.GetQueryInt("award", -1)
	ret, err := service.SubscribeRecord.GetList(userId, publisherId, pageNum, award)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *subscribeRecord) GetDetail(r *ghttp.Request) {
	orderNo := r.GetQueryString("orderNo")
	if orderNo == "" {
		s.FailJsonExit(r, "参数错误")
		return
	}
	userId := s.GetUserId(r)
	detail, err := service.SubscribeRecord.GetDetail(orderNo)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	if userId != detail.UserId {
		s.FailJsonExit(r, "无权访问")
		return
	}
	activityDetail, err := service.SubscribeActivity.GetSimpleDetail(detail.Aid)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	detail.ActivityType = activityDetail.ActivityType
	s.SusJsonExit(r, detail)
}

func (s *subscribeRecord) GetListByOrder(r *ghttp.Request) {
	pageNum := r.GetQueryInt("pageNum", 1)
	if pageNum <= 0 {
		pageNum = 1
	}
	status := r.GetQueryInt("status", -1)
	activityType := r.GetQueryInt("activityType")
	userId := s.GetUserId(r)
	publisherId := s.GetPublisherId(r)
	var aid = 0
	alias := r.GetQueryString("alias")
	if alias != "" {
		ainfo := service.SubscribeActivity.GetByAlias(alias)
		aid = ainfo.Id
	}
	ret, err := service.SubscribeRecord.GetListByOrder(userId, "", pageNum, status, publisherId, activityType, aid)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *subscribeRecord) GetDetailByOrder(r *ghttp.Request) {
	orderNo := r.GetQueryString("orderNo")
	if orderNo == "" {
		s.FailJsonExit(r, "参数错误")
		return
	}
	userId := s.GetUserId(r)
	publisherId := s.GetPublisherId(r)
	ret, err := service.SubscribeRecord.GetListByOrder(userId, orderNo, 1, -1, publisherId, 0, 0)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	if len(ret.List) == 0 {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret.List[0])
}

func (s *subscribeRecord) CreateOrder(r *ghttp.Request) {
	var req model.CreateOrderReq
	err := r.Parse(&req)
	if err != nil {
		s.FailJsonExit(r, "参数错误")
		return
	}
	if req.OrderNo == "" || req.SuccessRedirectUrl == "" || req.ExitRedirectUrl == "" {
		s.FailJsonExit(r, "参数错误")
		return
	}
	userId := s.GetUserId(r)
	ret, err := service.SubscribeRecord.CreateOrder(userId, r.GetClientIp(), req.OrderNo, req.SuccessRedirectUrl, req.ExitRedirectUrl, s.GetPublisherId(r), s.GetAppid(r))
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret.AppOrderNo)
}
