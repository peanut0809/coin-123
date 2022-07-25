package api

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/api"
	"github.com/gogf/gf/net/ghttp"
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
	s.SusJsonExit(r, ret)
}
