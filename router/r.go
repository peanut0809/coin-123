package router

import (
	"peanut-coin123/api/admin"
	frontApi "peanut-coin123/api/front"

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
		group.GET("/coin/ieo/off/items", admin.Admin.CoinIeoOffList)     // ieo、下架公告
		group.POST("/coin/core/create", admin.Admin.CreateCoreContext)   // 核心公告
		group.GET("/coin/core/items", admin.Admin.CoinCoreList)          // 核心公告
		group.POST("/coin/rename/create", admin.Admin.CreateCoinRename)  // 改名公告
		group.GET("/coin/rename/items", admin.Admin.CoinRenameList)      // 改名公告
	})
	s.Group("/front", func(group *ghttp.RouterGroup) {
		group.GET("/coin/items", frontApi.Front.CoinList)              // 现货交易对列表
		group.GET("/coin/ieo/off/items", frontApi.Front.CoinList)      // 币安ieo列表
		group.GET("/coin/rename/items", frontApi.Front.CoinRenameList) // 改名公告
		group.GET("/coin/core/items", frontApi.Front.CoinCoreList)     // 核心公告

	})
	return s
}
