package api

import (
	"fmt"
	"meta_launchpad/model"
	"meta_launchpad/service"

	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/api"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
)

type adminDropActivity struct {
	api.CommonBase
}

var AdminDropActivity = new(adminDropActivity)

func (s *adminDropActivity) Items(r *ghttp.Request) {
	var req *model.AirDropActivityItemReq
	if err := r.Parse(&req); err != nil {
		s.FailJsonExit(r, err.(gvalid.Error).FirstString())
		return
	}
	ret, err := service.AirDropActivity.Items(req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)

}

func (s *adminDropActivity) Item(r *ghttp.Request) {
	var req *model.AirDropActivityItemReq
	if err := r.Parse(&req); err != nil {
		s.FailJsonExit(r, err.(gvalid.Error).FirstString())
		return
	}
	ret, err := service.AirDropActivity.Item(req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *adminDropActivity) AirDrop(r *ghttp.Request) {

	var req *model.AirDropActivityReq
	if err := r.Parse(&req); err != nil {
		s.FailJsonExit(r, err.(gvalid.Error).FirstString())
		return
	}

	locKey := "meta_launchpad:admin:air:drop"
	re, e := g.Redis().Do("SET", locKey, 1, "ex", 60, "nx")
	if fmt.Sprintf("%v", re) == "OK" && e == nil {
		defer g.Redis().Do("DEL", locKey)
		ret, err := service.AirDropActivity.AirDrop(req)
		if err != nil {
			s.FailJsonExit(r, err.Error())
			return
		}
		s.SusJsonExit(r, ret)
		return
	} else {
		s.FailJsonExit(r, "操作太快了，请稍后重试")
		return
	}
}
