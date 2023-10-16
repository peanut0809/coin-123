package api

import (
	"meta_launchpad/service"

	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/api"
	"github.com/gogf/gf/net/ghttp"
)

type activity struct {
	api.CommonBase
}

var Activity = new(activity)

func (s *activity) ListBySearch(r *ghttp.Request) {
	pageNum := r.GetInt("pageNum", 1)
	pageSize := r.GetInt("pageSize", 20)
	searchVal := r.GetString("searchVal")
	ret, err := service.Activity.List(nil, pageNum, pageSize, "", "", 0, "", searchVal, "", 0)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *activity) GetPriceRank(r *ghttp.Request) {
	rankValue := r.GetQueryInt("rankValue")
	pageNum := r.GetQueryInt("pageNum", 1)
	pageSize := r.GetQueryInt("pageSize", 20)
	publisherId := r.GetQueryString("publisherId")
	searchVal := r.GetQueryString("searchVal")
	ret, err := service.Activity.GetCreatorRank(rankValue, pageNum, pageSize, publisherId, searchVal)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *activity) ListByClient(r *ghttp.Request) {
	pageNum := r.GetQueryInt("pageNum", 1)
	pageSize := r.GetQueryInt("pageSize", 20)
	activityType := r.GetQueryInt("activityType")
	searchVal := r.GetQueryString("searchVal")
	publisherId := s.GetPublisherId(r)
	if publisherId == "" {
		publisherId = r.GetQueryString("publisherId")
	}
	if publisherId == "" {
		s.FailJsonExit(r, "缺少发行商ID")
		return
	}
	ret, err := service.Activity.List(nil, pageNum, pageSize, "", "", activityType, "", searchVal, publisherId, 0)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}
func (s *activity) AllListByClient(r *ghttp.Request) {
	pageNum := r.GetQueryInt("pageNum", 1)
	pageSize := r.GetQueryInt("pageSize", 20)
	activityType := r.GetQueryInt("activityType")
	searchVal := r.GetQueryString("searchVal")
	publisherId := s.GetPublisherId(r)
	// if publisherId == "" {
	// 	publisherId = r.GetQueryString("publisherId")
	// }
	ret, err := service.Activity.AllList(nil, pageNum, pageSize, "", "", activityType, "", searchVal, publisherId, 0)
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
	ret, err := service.Activity.List(nil, pageNum, pageSize, startTime, endTime, activityType, status, searchVal, publisherId, -1)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}
