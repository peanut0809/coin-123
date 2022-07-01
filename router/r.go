package router

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/middleware"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func InitRouter() *ghttp.Server {
	s := g.Server()
	//跨域处理
	s.Use(middleware.CORS)
	// s.Group("/open", func(group *ghttp.RouterGroup) {
	// 	group.GET("/detail", api.AssetApi.GetDetail)
	// })
	// s.Group("/api", func(group *ghttp.RouterGroup) {
	// 	group.Middleware(middleware.Auth)
	// 	group.GET("/relation", api.AssetApi.GetRelation)
	// })
	// s.Group("/assets", func(group *ghttp.RouterGroup) {
	// 	group.Group("/product", func(group *ghttp.RouterGroup) {
	// 		group.POST("/list", api.AssetProduct.List)
	// 		group.GET("/detail", api.AssetProduct.Get)
	// 		group.GET("/asset", api.AssetProduct.GetAsset)
	// 	})
	// 	group.Group("/category", func(group *ghttp.RouterGroup) {
	// 		group.POST("/list", api.AssetProductCategory.ListAll)
	// 	})
	// })
	return s
}
