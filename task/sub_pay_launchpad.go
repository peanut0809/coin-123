package task

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/client"
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"meta_launchpad/provider"
	"meta_launchpad/service"
	"time"
)

type Extra struct {
	FromUserId string `json:"fromUserId"`
	ToUserId   string `json:"toUserId"`
	OrderNo    string `json:"orderNo"`
}

type RunSubLaunchpadPayData struct {
	Extra      string `json:"extra"`
	PayType    string `json:"pay_type"`
	AppOrderNo string `json:"app_order_no"`
}

func RunSubLaunchpadPayTask() {
	queue, err := client.GetQueue(client.QueueConfig{
		QueueName:  "payment.success.launchpad_pay",
		Exchange:   "payment.success.launchpad_pay",
		RoutingKey: "",
		Kind:       "fanout",
		MqUrl:      g.Cfg().GetString("rabbitmq.default.link"),
	})
	if err != nil {
		g.Log().Error("启动mq失败：" + err.Error())
		return
	}
	queue.RunConsume(func(msg []byte) error {
		g.Log().Info("RunSubLaunchpadPayTask receive:%s", gconv.String(msg))
		var data RunSubLaunchpadPayData
		err = json.Unmarshal(msg, &data)
		if err != nil {
			g.Log().Errorf("RunSubLaunchpadPayTask json.Unmarshal err:%v", err)
		} else {
			extra := Extra{}
			err = json.Unmarshal([]byte(data.Extra), &extra)
			if err != nil {
				g.Log().Errorf("RunSubLaunchpadPayTask json.Unmarshal err:%v", err)
			} else {
				//extra["orderNo"]
				err = service.SubscribeRecord.DoPaid(data.PayType, extra.OrderNo, data.AppOrderNo)
				if err != nil {
					g.Log().Errorf("RunSubLaunchpadPayTask err:%v", err)
				} else {
					subRecord, err := service.SubscribeRecord.GetSimpleDetail(extra.OrderNo)
					if err != nil {
						g.Log().Errorf("RunSubLaunchpadPayTask err:%v", err)
					} else {
						activityInfo, err := service.SubscribeActivity.GetSimpleDetail(subRecord.Aid)
						if err != nil {
							g.Log().Errorf("RunSubLaunchpadPayTask err:%v", err)
						} else {
							//发送资产到背包,重试三次
							publishSuccess := false
							for i := 0; i < 3; i++ {
								err = provider.Asset.PublishAssetWithTemplateId(&map[string]interface{}{
									"appId":      activityInfo.AppId,
									"templateId": activityInfo.TemplateId,
									"num":        subRecord.AwardNum,
									"userId":     subRecord.UserId,
									"optType":    "LAUNCHPAD",
									"optRemark":  "元初发射台发放资产",
									"nfrTime":    activityInfo.NfrSec,
								})
								if err != nil {
									g.Log().Errorf("RunSubLaunchpadPayTask err:%v 重试次数：%d", err, i)
									time.Sleep(time.Second)
									continue
								} else {
									publishSuccess = true
									break
								}
							}
							//更新发送资产的状态
							if publishSuccess {
								err = service.SubscribeRecord.UpdatePublishAsset(extra.OrderNo)
								if err != nil {
									g.Log().Errorf("RunSubLaunchpadPayTask err:%v", err)
									return nil
								}
								//通知钱包
								if data.PayType == "wallet_pay" {
									service.NoticeWallet(service.NoticeWalletReq{
										FromUserId: subRecord.UserId,
										ToUserId:   "B",
										TotalFee:   subRecord.SumPrice,
										OrderInfo:  gconv.String(msg),
									})
								}
							}
						}
					}
				}
			}
		}
		return nil
	})
}
