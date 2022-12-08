package api

import (
	"meta_launchpad/model"
	"meta_launchpad/service"

	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/api"
	"github.com/gogf/gf/net/ghttp"
)

type subscribeShare struct {
	api.CommonBase
}

var SubscribeShare = new(subscribeShare)

func (s *subscribeShare) UploadSubscribeShare(r *ghttp.Request) {
	var req model.SubscribeShareUpload
	err := r.Parse(&req)
	if err != nil {
		s.FailJsonExit(r, "参数错误")
		return
	}
	if req.Alias == "" {
		s.FailJsonExit(r, "参数错误")
		return
	}
	if req.PublisherId == "" {
		s.FailJsonExit(r, "参数错误")
		return
	}

	userId := s.GetUserId(r)
	err = service.SubscribeShare.UploadSubscrubeShare(req, userId)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, "success")
}
