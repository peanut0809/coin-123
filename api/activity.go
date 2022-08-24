package api

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/api"
	"github.com/gogf/gf/net/ghttp"
	"meta_launchpad/service"
)

type activity struct {
	api.CommonBase
}

var Activity = new(activity)

func (s *activity) Detail(r *ghttp.Request) {

}

func (s *activity) ListByClient(r *ghttp.Request) {
	pageNum := r.GetQueryInt("pageNum", 1)
	pageSize := r.GetQueryInt("pageSize", 20)
	publisherId := s.GetPublisherId(r)
	if publisherId == "" {
		publisherId = r.GetQueryString("publisherId")
	}
	if publisherId == "" {
		s.FailJsonExit(r, "缺少发行商ID")
		return
	}
	ret, err := service.Activity.List(nil, pageNum, pageSize, "", "", 0, "", "", publisherId)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *activity) List(r *ghttp.Request) {
	pageNum := r.GetQueryInt("pageNum", 1)
	pageSize := r.GetQueryInt("pageSize", 20)
	startTime := r.GetQueryString("startTime")
	endTime := r.GetQueryString("endTime")
	status := r.GetQueryString("status")
	activityType := r.GetQueryInt("activityType")
	publisherId := s.GetPublisherId(r)
	searchVal := r.GetQueryString("searchVal")
	ret, err := service.Activity.List(nil, pageNum, pageSize, startTime, endTime, activityType, status, searchVal, publisherId)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}
