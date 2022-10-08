package task

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/client"
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	"context"
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"meta_launchpad/model"
	"meta_launchpad/provider"
	"meta_launchpad/service"
	"time"
)

func DropTask() {
	queueName := "drop.asset"
	queue, err := client.GetQueue(client.QueueConfig{
		QueueName:  queueName,
		Exchange:   queueName,
		RoutingKey: "",
		Kind:       "fanout",
		MqUrl:      g.Cfg().GetString("rabbitmq.default.link"),
	})
	if err != nil {
		g.Log().Error("启动mq失败：" + err.Error())
		return
	}
	queue.RunConsume(func(msg []byte) error {
		g.Log().Info("DropTask receive:%s", gconv.String(msg))
		var data model.DropDetailRecord
		err = json.Unmarshal(msg, &data)
		if err != nil {
			return nil
		}
		assets, e := provider.Asset.GetCanUsedAssetsByTemplate(&map[string]interface{}{
			"appId":      data.AppId,
			"templateId": data.TemplateId,
			"num":        1,
		})
		if len(assets) == 0 {
			service.Drop.UpdateDetailRecordStatus(data.Id, 2, "", "库存不足")
			return nil
		}
		if e != nil {
			service.Drop.UpdateDetailRecordStatus(data.Id, 2, "", e.Error())
			return nil
		}
		tokenId := assets[0].TokenId
		list := make([]map[string]interface{}, 0)
		item := make(map[string]interface{})
		item["userId"] = data.UserId
		item["appId"] = data.AppId
		item["tokenId"] = tokenId
		item["nfrTime"] = time.Now().Add(time.Second * time.Duration(data.NfrSec)).Format("2006-01-02 15:04:05")
		list = append(list, item)
		params := &map[string]interface{}{
			"list": list,
			"opt":  map[string]interface{}{"optUserId": data.UserId, "optType": "DROP", "optRemark": "空投"},
		}
		_, err = utils.SendJsonRpc(context.Background(), "knapsack", "AssetKnapsack.Add", params)
		if err != nil {
			service.Drop.UpdateDetailRecordStatus(data.Id, 2, "", "发放背包失败")
			return nil
		}
		service.Drop.UpdateDetailRecordStatus(data.Id, 1, tokenId, "空投成功")
		return nil
	})
}
