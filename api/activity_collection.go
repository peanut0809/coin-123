package api

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/api"
	"github.com/gogf/gf/net/ghttp"
	"meta_launchpad/model"
	"meta_launchpad/service"
)

type activityCollection struct {
	api.CommonBase
}

var ActivityCollection = new(activityCollection)

func (s *activityCollection) ListByClient(r *ghttp.Request) {
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

	ret, err := service.ActivityCollection.ListByClient(0, publisherId, pageNum, pageSize)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *activityCollection) ListByDetail(r *ghttp.Request) {
	id := r.GetQueryInt("id")
	publisherId := s.GetPublisherId(r)
	if publisherId == "" {
		publisherId = r.GetQueryString("publisherId")
	}
	if publisherId == "" {
		s.FailJsonExit(r, "缺少发行商ID")
		return
	}
	ret, err := service.ActivityCollection.ListByClient(id, publisherId, 1, 1)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	if len(ret.List) == 0 {
		s.FailJsonExit(r, "活动不存在")
		return
	}
	ids, e := service.ActivityCollectionContent.GetActivityIds(ret.List[0].Id)
	if e != nil {
		s.FailJsonExit(r, e.Error())
		return
	}
	listInfo, e := service.Activity.List(ids, 1, 100, "", "", 0, "", "", publisherId)
	if e != nil {
		s.FailJsonExit(r, e.Error())
		return
	}
	response := model.ClientActivityCollectionDetail{}
	response.ActivityCollectionFull = ret.List[0]
	response.List = listInfo.List
	s.SusJsonExit(r, response)
}

func (s *activityCollection) List(r *ghttp.Request) {
	createStartTime := r.GetQueryString("createStartTime")
	createEndTime := r.GetQueryString("createEndTime")
	showStartTime := r.GetQueryString("showStartTime")
	showEndTime := r.GetQueryString("showEndTime")
	searchVal := r.GetQueryString("searchVal")
	pageNum := r.GetQueryInt("pageNum", 1)
	pageSize := r.GetQueryInt("pageSize", 20)
	publisherId := s.GetPublisherId(r)
	ret, err := service.ActivityCollection.List(publisherId, pageNum, createStartTime, createEndTime, showStartTime, showEndTime, searchVal, pageSize)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *activityCollection) Detail(r *ghttp.Request) {
	publisherId := s.GetPublisherId(r)
	id := r.GetQueryInt("id")
	ret, err := service.ActivityCollection.Detail(id, publisherId)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *activityCollection) Delete(r *ghttp.Request) {
	publisherId := s.GetPublisherId(r)
	id := r.GetInt("id")
	err := service.ActivityCollection.Delete(id, publisherId)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r)
}

func (s *activityCollection) Create(r *ghttp.Request) {
	var req model.CreateActivityCollectionReq
	err := r.Parse(&req)
	if err != nil {
		s.FailJsonExit(r, "参数错误")
		return
	}
	if req.Name == "" || req.Remark == "" || req.Intro == "" || req.ShowStartTime == nil || req.ShowEndTime == nil || req.Cover == "" {
		s.FailJsonExit(r, "参数错误")
		return
	}
	if req.ShowStartTime.Unix() > req.ShowEndTime.Unix() {
		s.FailJsonExit(r, "开始时间不能大于结束时间")
		return
	}
	if len(req.Activities) == 0 {
		s.FailJsonExit(r, "未选择营销活动")
		return
	}
	req.PublisherId = s.GetPublisherId(r)
	if req.Id == 0 {
		err = service.ActivityCollection.Create(req)
		if err != nil {
			s.FailJsonExit(r, err.Error())
			return
		}
		s.SusJsonExit(r)
		return
	}
	//更新
	err = service.ActivityCollection.Update(req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r)
	return
}
