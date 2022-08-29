package main

import (
	_ "brq5j1d.gfanx.pro/meta_cloud/meta_service/app/assets/dao"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcfg"
	gfnacos "github.com/imloama/gf-nacos"
	"meta_launchpad/router"
	_ "meta_launchpad/rpc"
	"meta_launchpad/task"
	"time"
)

func main() {
	s := router.InitRouter()
	g.Cfg().SetPath("./")
	s.Plugin(&gfnacos.GfNacosPlugin{
		Listener: func(configStr string) {
			g.Log().Info("配置文件发生了更新！\n", configStr)
			gcfg.SetContent(configStr)

			go func() {
				time.Sleep(time.Second * 3)
				go task.RunSubTask()
				go task.RunSubPayTask()
				go task.RunSubLaunchpadPayTask()
				go task.RunSeckillOrderTask()
				go task.RunSeckillOrderPayTask()
			}()

			//err := service.Sms.SendSms("13720009841", "ecgDjLtq", "aIIbedlG", "4HZdAzLt", map[string]string{
			//	"goods": "中华网数藏印象·故宫·神秘",
			//	"time":  "04:05",
			//})
			//fmt.Println(err)
		},
	})

	go task.RunLuckyDrawTask()

	go task.CheckSubPayTimeout()

	go task.CheckSeckillOrderTimeoutTask()

	s.Run()
}
