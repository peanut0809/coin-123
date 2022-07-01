package router

import (
	apiMiddleware "meta_open_sdk/app/open/middleware"

	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/middleware"
	third "brq5j1d.gfanx.pro/meta_cloud/meta_service/app/third/middleware"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func InitRouter() *ghttp.Server {
	s := g.Server()
	//跨域处理
	s.Use(middleware.CORS)
	s.Group("/third", func(group *ghttp.RouterGroup) {
		group.Middleware(apiMiddleware.ApiCheck)
		group.Middleware(third.VerifySign)
		// group.POST("/order/create", api.OrderApi.CreateOrder)              //创建订单
		// group.POST("/token/verify", api.UserApi.VerifyGameToken)           //校验gameToken
		// group.POST("/user/create", api.UserApi.CreateUser)                 //创建游戏角色
		// group.POST("/user/crystal/add", api.UserApi.AddCrystal)            //消耗钻石或红钻获得元晶
		// group.POST("/user/inviter/info", api.UserApi.GetInviterInfo)       //查询用户邀请人信息
		// group.POST("/user/inviter/add", api.UserApi.UpdateInviterCode)     //添加用户上级邀请人
		// group.POST("/user/crystal/change", api.UserApi.AccountChange)      //元晶账户变动
		// group.POST("/user/crystal/balance", api.UserApi.GetCrystalBalance) //获取元晶的余额
	})
	s.Group("/third/open", func(group *ghttp.RouterGroup) {
		// group.POST("/user/loginfo/phone", api.UserApi.GetUserLoginfo) //通过手机号返回登录信息
	})
	return s
}
