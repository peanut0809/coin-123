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
