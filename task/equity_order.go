package task

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/client"
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"meta_launchpad/model"
	"meta_launchpad/service"
)

func RunEquityOrderTask() {
	queue, err := client.GetQueue(client.QueueConfig{
		QueueName:  "launchpad.equity",
		Exchange:   "launchpad.equity",
		RoutingKey: "",
		Kind:       "fanout",
		MqUrl:      g.Cfg().GetString("rabbitmq.default.link"),
	})
	if err != nil {
		g.Log().Error("启动mq失败：" + err.Error())
		return
	}
	queue.RunConsume(func(msg []byte) error {
		g.Log().Info("RunEquityOrderTask receive:%s", gconv.String(msg))
		var data model.EquityOrderReq
		err = json.Unmarshal(msg, &data)
		if err == nil {
			// 创建订单
			service.Equity.Create(data)
		} else {
			g.Log().Info("RunEquityOrderTask json.Unmarshal err:%v", err)
		}
		return nil
	})
}
