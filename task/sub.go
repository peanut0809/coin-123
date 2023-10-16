package task

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/client"
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"meta_launchpad/model"
	"meta_launchpad/service"
)

func RunSubTask() {
	queue, err := client.GetQueue(client.QueueConfig{
		QueueName:  "launchpad.sub",
		Exchange:   "launchpad.sub",
		RoutingKey: "",
		Kind:       "fanout",
		MqUrl:      g.Cfg().GetString("rabbitmq.default.link"),
	})
	if err != nil {
		g.Log().Error("启动mq失败：" + err.Error())
		return
	}
	queue.RunConsume(func(msg []byte) error {
		g.Log().Info("RunSubTask receive:%s", gconv.String(msg))
		var data model.DoSubReq
		err = json.Unmarshal(msg, &data)
		if err == nil {
			service.SubscribeActivity.DoSub(data)
		} else {
			g.Log().Info("RunSubTask json.Unmarshal err:%v", err)
		}
		return nil
	})
}
