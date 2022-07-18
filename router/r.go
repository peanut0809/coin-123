package router

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/middleware"
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/library"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"meta_launchpad/api"
)

func V(r *ghttp.Request) {
	loginFrom := r.GetCtxVar("loginFrom").String()
	if loginFrom != "LAUNCHPAD" {
		library.FailJsonCodeExit(r, gerror.NewCode(gcode.New(-403, "无权访问此接口", nil)))
		return
	}
	r.Middleware.Next()
}

func InitRouter() *ghttp.Server {
	s := g.Server()
	//跨域处理
	s.Use(middleware.CORS)
	s.Group("/open", func(group *ghttp.RouterGroup) {
		group.GET("/activity/list", api.SubscribeActivity.GetSubscribeActivityList)
	})
	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.Auth)
		group.GET("/activity/detail", api.SubscribeActivity.GetSubscribeActivityDetail)
		group.GET("/pay/info", api.SubscribeActivity.GetPayInfo)
		group.POST("/activity/sub", api.SubscribeActivity.SubActivity)
		group.GET("/activity/sub/result", api.SubscribeActivity.GetSubActivityResult)
		group.GET("/sub/list", api.SubscribeRecord.GetList)
		group.GET("/sub/detail", api.SubscribeRecord.GetDetail)
		group.GET("/sub/order/list", api.SubscribeRecord.GetListByOrder)
		group.GET("/sub/order/detail", api.SubscribeRecord.GetDetailByOrder)
		group.POST("/sub/order/create", api.SubscribeRecord.CreateOrder)
	})
	return s
}
