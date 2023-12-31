package task

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/client"
	"context"
	"encoding/json"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"meta_launchpad/model"
	"meta_launchpad/service"
)

type RunSubPayTaskData struct {
	Extra   string `json:"extra"`
	PayType string `json:"pay_type"`
}

func RunSubPayTask() {
	queue, err := client.GetQueue(client.QueueConfig{
		QueueName:  "payment.success.launchpad_ticket",
		Exchange:   "payment.success.launchpad_ticket",
		RoutingKey: "",
		Kind:       "fanout",
		MqUrl:      g.Cfg().GetString("rabbitmq.default.link"),
	})
	if err != nil {
		g.Log().Error("启动mq失败：" + err.Error())
		return
	}
	queue.RunConsume(func(msg []byte) error {
		g.Log().Info("RunSubPayTask receive:%s", gconv.String(msg))
		var data RunSubPayTaskData
		err = json.Unmarshal(msg, &data)
		if err != nil {
			g.Log().Errorf("RunSubPayTask json.Unmarshal err:%v", err)
		} else {
			var extra model.SubscribeRecordQueueData
			err = json.Unmarshal([]byte(data.Extra), &extra)
			if err != nil {
				g.Log().Errorf("RunSubPayTask json.Unmarshal err:%v", err)
			} else {
				extra.PayTicketMethod = data.PayType
				err = g.DB().Transaction(context.Background(), func(ctx context.Context, tx *gdb.TX) error {
					err = service.SubscribeRecord.CreateSubscribeRecord(tx, model.SubscribeRecord{
						ActivityType:    extra.ActivityType,
						Aid:             extra.Aid,
						Name:            extra.Name,
						Icon:            extra.Icon,
						UserId:          extra.UserId,
						BuyNum:          extra.BuyNum,
						OrderNo:         extra.OrderNo,
						PayOrderNo:      extra.OrderNo,
						SumPrice:        extra.SumPrice,
						SumUnitPrice:    extra.SumUnitPrice,
						TicketType:      extra.TicketType,
						PublisherId:     extra.PublisherId,
						PayEndTime:      extra.PayEndTime,
						PayTicketMethod: data.PayType,
					})
					return err
				})
				if err != nil {
					g.Log().Errorf("RunSubPayTask json.Unmarshal err:%v", err)
				} else {
					//通知钱包
					if data.PayType == "wallet_pay" && extra.SumUnitPrice != 0 {
						service.NoticeWallet(service.NoticeWalletReq{
							FromUserId: extra.UserId,
							ToUserId:   "B",
							TotalFee:   extra.SumUnitPrice,
							OrderInfo:  gconv.String(msg),
						})
					}
				}
			}
		}
		return nil
	})
}
