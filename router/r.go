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
		//group.GET("/activity/list", api.SubscribeActivity.GetSubscribeActivityList)
		group.GET("/seckill/activity/detail", api.SeckillActivity.GetDetail) //活动详情
		group.GET("/activity/award/record", api.SubscribeActivity.GetActivityAwardRecord)
		group.GET("/activity/detail", api.SubscribeActivity.GetSubscribeActivityDetail)
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

		//秒杀活动
		group.GET("/seckill/activity/detail", api.SeckillActivity.GetDetail)                         //活动详情
		group.POST("/seckill/activity/order/create", api.SeckillActivity.CreateOrder)                //下单
		group.POST("/seckill/activity/order/cancel", api.SeckillActivity.CancelOrder)                //取消订单
		group.GET("/seckill/activity/order/create/result", api.SeckillActivity.GetCreateOrderResult) //获取下单结果
		group.GET("/seckill/activity/order/list", api.SeckillActivity.GetOrderList)                  //订单列表
		group.GET("/seckill/activity/order/detail", api.SeckillActivity.GetOrderDetail)              //订单详情
	})

	s.Group("/admin", func(group *ghttp.RouterGroup) {
		group.Middleware(func(r *ghttp.Request) {
			//测试临时写个发行商
			r.SetCtxVar("publisherId", "TEST")
			r.Middleware.Next()
		})
		group.POST("/activity/create", api.AdminSubscribeActivity.Create)
		group.GET("/activity/list", api.AdminSubscribeActivity.List)
		group.GET("/activity/detail", api.AdminSubscribeActivity.Detail)
		group.POST("/activity/delete", api.AdminSubscribeActivity.Delete)
		group.POST("/activity/disable", api.AdminSubscribeActivity.Disable)
		group.GET("/activity/sub/record", api.AdminSubscribeActivity.GetSubRecords)

		//group.POST("/activity/sub/record", api.AdminSubscribeActivity.GetSubRecords)

		//秒杀
		group.POST("/seckill/activity/create", api.AdminSeckillActivity.Create)
		group.GET("/seckill/activity/detail", api.AdminSeckillActivity.Detail)

		//banner
		group.GET("/banner/list", api.Banner.GetList)
		group.POST("/banner/create", api.Banner.Create)
		group.DELETE("/banner/delete", api.Banner.Delete)
		group.PUT("/banner/stateEdit", api.Banner.StateEdit)
	})
	return s
}
