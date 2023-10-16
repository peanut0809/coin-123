package api

import (
	"github.com/gogf/gf/net/ghttp"
)

type commonApi struct {
}

var CommonApi = new(commonApi)

func (s *commonApi) Items(r *ghttp.Request) {
	return
}
