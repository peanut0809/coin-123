package router

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/middleware"
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/library"
	"fmt"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"meta_launchpad/api"
	"meta_launchpad/cache"
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
		group.GET("/temp/del", func(r *ghttp.Request) {
			//检查超时行为
			userId := r.GetQueryString("userId")
			aType := r.GetQueryInt("aType")
			fmt.Println("=======", fmt.Sprintf(cache.SUB_PAY_TIMEOUT, userId, aType))
			g.Redis().Do("DEL", fmt.Sprintf(cache.SUB_PAY_TIMEOUT, userId, aType))

			gv, _ := g.Redis().DoVar("GET", fmt.Sprintf(cache.SUB_PAY_TIMEOUT, userId, aType))
			fmt.Println("=======", gv.IsEmpty())
		})
	})
	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.Auth)
		group.GET("/activity/detail", api.SubscribeActivity.GetSubscribeActivityDetail)
		group.GET("/activity/award/record", api.SubscribeActivity.GetActivityAwardRecord)
		group.GET("/pay/info", api.SubscribeActivity.GetPayInfo)
		group.POST("/activity/sub", api.SubscribeActivity.SubActivity)
		group.GET("/activity/sub/result", api.SubscribeActivity.GetSubActivityResult)
		group.GET("/activity/award/result", api.SubscribeActivity.GetSubActivityAwardResult)
		group.GET("/sub/list", api.SubscribeRecord.GetList)
		group.GET("/sub/detail", api.SubscribeRecord.GetDetail)
		group.GET("/sub/order/list", api.SubscribeRecord.GetListByOrder)
		group.GET("/sub/order/detail", api.SubscribeRecord.GetDetailByOrder)
		group.POST("/sub/order/create", api.SubscribeRecord.CreateOrder)
	})
	return s
}
