package api

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/library"
	"context"
	"encoding/json"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/net/ghttp"
)

// GetPublisherByToken 通过token获取发行商ID
func GetPublisherByToken(r *ghttp.Request) {
	token := r.Header.Get("x-token")
	params := &map[string]interface{}{
		"token": token,
	}
	result, err := utils.SendJsonRpc(context.Background(), "developer", "Publisher.GetPublisherByToken", params)
	s := struct {
		PublisherId string `json:"publisherId"`
		Code        int    `json:"code"`
	}{}
	marshal, err := json.Marshal(result)
	if err != nil {
		return
	}
	err = json.Unmarshal(marshal, &s)
	if err != nil {
		return
	}
	if s.Code == -401 {
		library.FailJsonCodeExit(r, gerror.NewCode(gcode.New(s.Code, "没有此发行商", nil)))
	} else if s.Code == 404 {
		library.FailJsonCodeExit(r, gerror.NewCode(gcode.New(s.Code, "token无效", nil)))
	}
	if s.PublisherId == "" {
		library.FailJsonCodeExit(r, gerror.NewCode(gcode.New(s.Code, "发行商错误", nil)))
	}
	r.SetCtxVar("publisherId", s.PublisherId)
	r.Middleware.Next()
}
