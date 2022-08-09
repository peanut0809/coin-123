package service

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/client"
	"github.com/gogf/gf/frame/g"
)

type NoticeWalletReq struct {
	FromUserId string `json:"fromUserId"`
	ToUserId   string `json:"toUserId"`
	Gas        int    `json:"gas"`
	TotalFee   int    `json:"totalFee"`
	AppId      string `json:"appId"`
	TokenId    string `json:"tokenId"`
	OrderInfo  string `json:"orderInfo"`
}

func NoticeWallet(in NoticeWalletReq) {
	queueName := "market.notification.wallet"
	mqClient, err := client.GetQueue(client.QueueConfig{
		QueueName:  queueName,
		Exchange:   queueName,
		RoutingKey: "",
		Kind:       "fanout",
		MqUrl:      g.Cfg().GetString("rabbitmq.default.link"),
	})
	if err != nil {
		g.Log().Errorf("NoticeWallet err:%v，data:%+v", err, in)
		return
	}
	err = mqClient.Publish(in)
	if err != nil {
		g.Log().Errorf("NoticeWallet err:%v，data:%+v", err, in)
		return
	}
}
