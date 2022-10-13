package api

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/api"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"meta_launchpad/model"
	"meta_launchpad/service"
)

type synthetic struct {
	api.CommonBase
}

var Synthetic = new(synthetic)

func (s *synthetic) Create(r *ghttp.Request) {
	var req model.SyntheticActivityReq
	err := r.Parse(&req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	req.PublisherId = s.GetPublisherId(r)
	if req.StartTime == nil || req.EndTime == nil || req.Cover == "" || req.Name == "" || req.AppId == "" || req.TemplateId == "" || req.Sum <= 0 || req.OutNum <= 0 || len(req.ConditionArr) == 0 {
		s.FailJsonExit(r, "参数错误")
		return
	}
	if req.StartTime.Unix() >= req.EndTime.Unix() {
		s.FailJsonExit(r, "结束时间必须开始时间")
		return
	}
	conditionStr := gconv.String(req.ConditionArr)
	req.SyntheticActivity.RemainNum = req.SyntheticActivity.Sum
	req.SyntheticActivity.Condition = &conditionStr
	if req.Id == 0 {
		err = service.Synthetic.Create(req.SyntheticActivity)
		if err != nil {
			s.FailJsonExit(r, err.Error())
			return
		}
		s.SusJsonExit(r)
		return
	}
	err = service.Synthetic.Update(req.SyntheticActivity)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r)
	return
}

func (s *synthetic) List(r *ghttp.Request) {
	publisherId := s.GetPublisherId(r)
	pageNum := r.GetQueryInt("pageNum", 1)
	pageSize := r.GetQueryInt("pageSize", 1)
	startTimeBegin := r.GetQueryString("startTimeBegin")
	startTimeEnd := r.GetQueryString("startTimeEnd")
	endTimeBegin := r.GetQueryString("endTimeBegin")
	endTimeEnd := r.GetQueryString("endTimeEnd")
	status := r.GetQueryString("status")
	searchVal := r.GetQueryString("searchVal")
	ret, err := service.Synthetic.List(publisherId, pageNum, pageSize, startTimeBegin, startTimeEnd, endTimeBegin, endTimeEnd, status, searchVal, "id DESC", 0)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
	return
}

func (s *synthetic) ClientList(r *ghttp.Request) {
	publisherId := s.GetPublisherId(r)
	pageNum := r.GetQueryInt("pageNum", 1)
	pageSize := r.GetQueryInt("pageSize", 1)
	status := r.GetQueryString("status")
	ret, err := service.Synthetic.List(publisherId, pageNum, pageSize, "", "", "", "", status, "", "start_time DESC", 1)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
	return
}

func (s *synthetic) Detail(r *ghttp.Request) {
	id := r.GetQueryInt("id")
	ret, err := service.Synthetic.Detail(id)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
	return
}

func (s *synthetic) Open(r *ghttp.Request) {
	id := r.GetInt("id")
	open := r.GetInt("open")
	ret, err := service.Synthetic.Open(id, open)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
	return
}

func (s *synthetic) ClientDetail(r *ghttp.Request) {
	id := r.GetInt("id")
	ret, err := service.Synthetic.ClientDetail(id)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}
