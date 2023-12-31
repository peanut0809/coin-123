package router

import (
	"fmt"
	"meta_launchpad/api"
	"meta_launchpad/cache"

	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/middleware"
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/library"
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/v2/common/admin"
	sysApi "brq5j1d.gfanx.pro/meta_cloud/meta_common/v2/common/api"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
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
		group.GET("/sass/activity/alllist", api.Activity.AllListByClient)
		//身价排行榜
		group.GET("/sass/price/rank", api.Activity.GetPriceRank)

		//市场搜索
		group.POST("/sass/activity/search", api.Activity.ListBySearch)
		//活动合集
		group.GET("/sass/activity/collection/list", api.ActivityCollection.ListByClient)
		//活动合集详情
		group.GET("/sass/activity/collection/detail", api.ActivityCollection.ListByDetail)

		// 权益购
		group.GET("/activity/equity/list", api.Equity.List) // 活动列表
		group.GET("/activity/equity/info", api.Equity.Info) // 活动详情

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
		group.POST("/sub/share/upload", api.SubscribeShare.UploadSubscribeShare)

		//秒杀活动
		group.GET("/seckill/activity/detail", api.SeckillActivity.GetDetail)                         //活动详情
		group.POST("/seckill/activity/order/create", api.SeckillActivity.CreateOrder)                //下单
		group.POST("/seckill/activity/order/cancel", api.SeckillActivity.CancelOrder)                //取消订单
		group.GET("/seckill/activity/order/create/result", api.SeckillActivity.GetCreateOrderResult) //获取下单结果
		group.GET("/seckill/activity/order/list", api.SeckillActivity.GetOrderList)                  //订单列表
		group.GET("/seckill/activity/order/detail", api.SeckillActivity.GetOrderDetail)              //订单详情

		// 免费版后台首页接口
		group.GET("/frontPage/freeTransactionSlip", api.FrontPage.FreeTransactionSlip) // 交易数，交易总数
		group.GET("/frontPage/freeVolumeOfTrade", api.FrontPage.FreeVolumeOfTrade)     // 近期支付数
		group.GET("/frontPage/freeTransactionsNum", api.FrontPage.FreeTransactionsNum) // 支付笔数
		group.GET("/frontPage/freePayers", api.FrontPage.FreePayers)                   // 支付人数
		group.GET("/frontPage/freeTurnover", api.FrontPage.FreeTurnover)               // 交易额

		//合成
		group.GET("/synthetic/list", api.Synthetic.ClientList)
		group.GET("/synthetic/detail", api.Synthetic.ClientDetail)
		group.POST("/synthetic/do", api.Synthetic.DoSynthetic)
		group.GET("/synthetic/do/result", api.Synthetic.GetDoSyntheticResult)
		group.GET("/synthetic/record/list", api.Synthetic.GetRecordList)
		group.GET("/synthetic/record/detail", api.Synthetic.GetRecordDetail)
		// 权益活动
		group.POST("/activity/equity/order/create", api.Equity.CreateOrder)                // 活动下单
		group.POST("/activity/equity/order/cancel", api.Equity.CancelOrder)                // 取消订单
		group.GET("/activity/equity/order/create/result", api.Equity.GetCreateOrderResult) // 下单结果
		group.GET("/activity/equity/order/list", api.Equity.GetOrderList)                  // 订单列表
		group.GET("/activity/equity/order/detail", api.Equity.GetOrderDetail)              // 订单详情
		group.GET("/activity/equity/order/buynum", api.Equity.CanBuyNum)                   // 可购买数量
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
		group.GET("/synthetic/record", api.Synthetic.GetSyntheticRecord)

		// 权益活动后台@段昀彤
		group.GET("/equity/assets/count", api.AdminEquity.AssetsCount)          //获取可用资产数
		group.POST("/equity/activity/create", api.AdminEquity.Create)           //白名单活动创建
		group.GET("/equity/activity/invalid", api.AdminEquity.Invalid)          //下架更新
		group.POST("/equity/activity/user/import", api.AdminEquity.Import)      //用户导入
		group.GET("/equity/activity/item", api.AdminEquity.Item)                //白名单活动创建
		group.GET("/equity/activity/user/items", api.AdminEquity.UserItems)     //权益活动关联用户
		group.GET("/equity/activity/items", api.AdminEquity.EquityItems)        //权益活动列表
		group.GET("/equity/activity/put/items", api.AdminEquity.EquityPutItems) //获取已上架活动列表
		//  rpc  定时任务处理活动下架状态
		group.GET("/equity/order/items", api.AdminEquity.OrderItems)    //订单查询
		group.POST("/equity/order/export", api.AdminEquity.OrderExport) //订单导出

	})

	s.Group("/admin", func(group *ghttp.RouterGroup) {

		group.Middleware(sysApi.Auth.VerifyToken)
		group.Middleware(admin.Ctx, admin.Auth)

		group.Hook("/*", ghttp.HookAfterOutput, admin.OptionLog)
		// 空投配置集合
		group.POST("/air/items", api.AdminDropActivity.Items)
		// 空投配置明细集合
		group.POST("/air/item", api.AdminDropActivity.Item)
		// 发起空投
		group.POST("/air", api.AdminDropActivity.AirDrop)

	})

	return s
}
