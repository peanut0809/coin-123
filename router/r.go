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
		//banner
		group.GET("/banner/getFrontList", api.Banner.GetFrontList)
		group.GET("/banner/getRichText", api.Banner.GetRichText)

		//最新上线
		group.GET("/sass/activity/list", api.Activity.ListByClient)
		//市场搜索
		group.POST("/sass/activity/search", api.Activity.ListBySearch)
		//活动合集
		group.GET("/sass/activity/collection/list", api.ActivityCollection.ListByClient)
		//活动合集详情
		group.GET("/sass/activity/collection/detail", api.ActivityCollection.ListByDetail)

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

		//合成
		group.GET("/synthetic/list", api.Synthetic.ClientList)
		group.GET("/synthetic/detail", api.Synthetic.ClientDetail)
	})

	s.Group("/admin", func(group *ghttp.RouterGroup) {
		group.Middleware(api.GetPublisherByToken)
		group.POST("/sub/activity/create", api.AdminSubscribeActivity.Create)
		group.GET("/sub/activity/list", api.AdminSubscribeActivity.List)
		group.GET("/sub/activity/detail", api.AdminSubscribeActivity.Detail)
		group.POST("/sub/activity/delete", api.AdminSubscribeActivity.Delete)
		group.POST("/sub/activity/disable", api.AdminSubscribeActivity.Disable)
		group.GET("/sub/activity/record", api.AdminSubscribeActivity.GetSubRecords)

		//活动
		group.GET("/activity/list", api.Activity.List)

		//活动合集
		group.POST("/activity/collection/create", api.ActivityCollection.Create)
		group.GET("/activity/collection/detail", api.ActivityCollection.Detail)
		group.GET("/activity/collection/list", api.ActivityCollection.List)
		group.POST("/activity/collection/delete", api.ActivityCollection.Delete)

		//秒杀
		group.POST("/seckill/activity/create", api.AdminSeckillActivity.Create)
		group.GET("/seckill/activity/detail", api.AdminSeckillActivity.Detail)
		group.POST("/seckill/activity/disable", api.AdminSeckillActivity.Disable)
		group.GET("/seckill/activity/list", api.AdminSeckillActivity.List)
		group.POST("/seckill/activity/delete", api.AdminSeckillActivity.Delete)
		group.GET("/seckill/activity/orders", api.AdminSeckillActivity.GetOrders)

		//空投
		group.POST("/drop/do", api.Drop.Create)
		group.GET("/drop/record/list", api.Drop.GetRecordList)
		group.GET("/drop/record/detail/list", api.Drop.GetDetailRecordList)

		//banner 后端接口
		group.GET("/banner/list", api.Banner.GetList)
		group.POST("/banner/create", api.Banner.Create)
		group.POST("/banner/delete", api.Banner.Delete)
		group.POST("/banner/stateEdit", api.Banner.StateEdit)

		// 后台首页接口
		group.GET("/frontPage/transactionSlip", api.FrontPage.TransactionSlip) // 交易数，交易总数
		group.GET("/frontPage/volumeOfTrade", api.FrontPage.VolumeOfTrade)     // 近期支付数
		group.GET("/frontPage/transactionsNum", api.FrontPage.TransactionsNum) // 支付笔数
		group.GET("/frontPage/payers", api.FrontPage.Payers)                   // 支付人数
		group.GET("/frontPage/turnover", api.FrontPage.Turnover)               // 交易额

		//合成
		group.POST("/synthetic/create", api.Synthetic.Create)
		group.GET("/synthetic/list", api.Synthetic.List)
		group.GET("/synthetic/detail", api.Synthetic.Detail)
		group.POST("/synthetic/open", api.Synthetic.Open)
		group.POST("/synthetic/delete", api.Synthetic.Delete)
	})
	return s
}
