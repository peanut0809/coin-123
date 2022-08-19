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

func (s *activityCollection) List(r *ghttp.Request) {

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
