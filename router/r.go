package router

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func InitRouter() *ghttp.Server {
	s := g.Server()

	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("Hello, World!")

	})
	return s
}
