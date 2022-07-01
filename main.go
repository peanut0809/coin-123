package main

import (
	"meta_open_sdk/router"

	_ "brq5j1d.gfanx.pro/meta_cloud/meta_service/app/assets/dao"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcfg"
	gfnacos "github.com/imloama/gf-nacos"
)

func main() {
	s := router.InitRouter()
	g.Cfg().SetPath("./")
	s.Plugin(&gfnacos.GfNacosPlugin{
		Listener: func(configStr string) {
			g.Log().Info("配置文件发生了更新！\n", configStr)
			gcfg.SetContent(configStr)
		},
	})
	s.Run()
}
