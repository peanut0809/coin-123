package router

import (
	"cccn-zxl-server/api"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func InitRouter() *ghttp.Server {
	s := g.Server()
	s.Group("/zxl", func(group *ghttp.RouterGroup) {
		group.POST("/order", api.ZxlApi.MarktingZxlOrder)          // zxl下单
		group.POST("/order/info", api.ZxlApi.MarktingZxlOrderInfo) // zxl下单
	})
	return s
}
