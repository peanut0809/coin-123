package router

import (
	"peanut-coin123/api/admin"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func InitRouter() *ghttp.Server {
	s := g.Server()
	s.Group("/common", func(group *ghttp.RouterGroup) {

	})
	s.Group("/admin", func(group *ghttp.RouterGroup) { // zxl下单
		group.POST("/coin/create", admin.Admin.CreateCoinItems)          // 现货交易对
		group.GET("/coin/items", admin.Admin.CoinList)                   // 现货交易对列表
		group.POST("/coin/ieo/off/create", admin.Admin.CreateIeoOffCoin) // ieo、下架公告
		group.GET("/coin/ieo/off/items", admin.Admin.CreateIeoOffCoin)   // ieo、下架公告
		group.POST("/coin/core/create", admin.Admin.CreateCoreContext)   // 核心公告
		group.POST("/coin/core/items", admin.Admin.CreateCoreContext)    // 核心公告
		group.POST("/coin/rename/create", admin.Admin.CreateCoinRename)  // 改名公告
		group.POST("/coin/rename/items", admin.Admin.CreateCoinRename)   // 改名公告

	})
	s.Group("/front", func(group *ghttp.RouterGroup) {
		// 币安上币集合
		// 币安/upbit交易排行榜top20
		// 行情推荐/费
		// 关于我们
	})
	return s
}
