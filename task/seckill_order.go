package task

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/client"
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"meta_launchpad/model"
	"meta_launchpad/service"
)

func RunSeckillOrderTask() {
	queue, err := client.GetQueue(client.QueueConfig{
		QueueName:  "launchpad.seckill",
		Exchange:   "launchpad.seckill",
		RoutingKey: "",
		Kind:       "fanout",
		MqUrl:      g.Cfg().GetString("rabbitmq.default.link"),
	})
	if err != nil {
		g.Log().Error("启动mq失败：" + err.Error())
		return
	}
	queue.RunConsume(func(msg []byte) error {
		g.Log().Info("RunSeckillOrderTask receive:%s", gconv.String(msg))
		var data model.DoBuyReq
		err = json.Unmarshal(msg, &data)
		if err == nil {
			service.SeckillActivity.DoBuy(data)
		} else {
			g.Log().Info("RunSeckillOrderTask json.Unmarshal err:%v", err)
		}
		return nil
	})
}
