package task

import (
	"encoding/json"
	"meta_launchpad/model"
	"meta_launchpad/service"

	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/client"
	"github.com/gogf/gf/frame/g"
)

func RunSyntheticTask() {
	queue, err := client.GetQueue(client.QueueConfig{
		QueueName:  "synthetic.do",
		Exchange:   "synthetic.do",
		RoutingKey: "",
		Kind:       "fanout",
		MqUrl:      g.Cfg().GetString("rabbitmq.default.link"),
	})
	if err != nil {
		g.Log().Error("启动mq失败：" + err.Error())
		return
	}
	queue.RunConsume(func(msg []byte) error {
		// g.Log().Info("RunSyntheticTask receive:%s", gconv.String(msg))
		var data model.SyntheticReq
		err = json.Unmarshal(msg, &data)
		if err == nil {
			service.Synthetic.Synthetic(data)
		} else {
			g.Log().Info("RunSyntheticTask json.Unmarshal err:%v", err)
		}
		return nil
	})
}
