package task

import (
	"encoding/json"
	"meta_launchpad/model"
	"meta_launchpad/provider"
	"meta_launchpad/service"
	"strings"

	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/client"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
)

func RunEquityOrderPayTask() {
	queue, err := client.GetQueue(client.QueueConfig{
		QueueName:  "payment.success.launchpad_equity",
		Exchange:   "payment.success.launchpad_equity",
		RoutingKey: "",
		Kind:       "fanout",
		MqUrl:      g.Cfg().GetString("rabbitmq.default.link"),
	})
	if err != nil {
		g.Log().Error("启动mq失败：" + err.Error())
		return
	}
	queue.RunConsume(func(msg []byte) error {
		g.Log().Info("RunEquityOrderPayTask receive:%s", gconv.String(msg))
		var data RunSubLaunchpadPayData
		err = json.Unmarshal(msg, &data)
		if err == nil {
			orderInfo, err := service.EquityOrder.GetInfoByOrderNo(data.AppOrderNo)
			if err != nil {
				g.Log().Errorf("RunEquityOrderPayTask err:%v", err)
				return nil
			}
			activityInfo, err := service.Equity.GetValidDetail(orderInfo.ActivityId)
			if err != nil {
				g.Log().Errorf("RunSubLaunchpadPayTask err:%v", err)
				return nil
			}
			tx, e := g.DB().Begin()
			if e != nil {
				g.Log().Info("RunSubLaunchpadPayTask err:%v", e)
			}
			r, e := tx.Exec("UPDATE equity_orders SET status = ?,pay_at = ?,pay_method = ? WHERE order_no = ?", model.PAID, gtime.Now(), data.PayType, data.AppOrderNo)
			if e != nil {
				g.Log().Info("RunSubLaunchpadPayTask err:%v", e)
				return nil
			}
			affectedNum, _ := r.RowsAffected()
			if affectedNum != 1 {
				err = tx.Rollback()
				return nil
			}
			// 发放资产
			publishSuccess := true
			err = provider.Asset.PublishAssetWithTemplateId(&map[string]interface{}{
				"appId":      activityInfo.AppId,
				"templateId": activityInfo.TemplateId,
				"num":        orderInfo.Num,
				"userId":     orderInfo.UserId,
				"optType":    "LAUNCHPAD",
				"optRemark":  "元初发射台权益发放资产",
				"nfrTime":    activityInfo.NfrSec,
			})
			if err != nil {
				if strings.Contains(err.Error(), "timeout") {
					publishSuccess = true
				} else {
					publishSuccess = false
				}
			}
			// 更新发送资产的状态
			if publishSuccess {
				_, e = g.DB().Exec("UPDATE equity_orders SET publish_status = 1 WHERE order_no = ?", data.AppOrderNo)
				if e != nil {
					g.Log().Info("RunSubLaunchpadPayTask err:%v", e)
					return nil
				}
				// 通知钱包
				if data.PayType == "wallet_pay" {
					service.NoticeWallet(service.NoticeWalletReq{
						FromUserId: orderInfo.UserId,
						ToUserId:   "B",
						TotalFee:   orderInfo.Price,
						OrderInfo:  gconv.String(msg),
					})
				}
			} else {
				_, e = g.DB().Exec("UPDATE equity_orders SET publish_status = 2 WHERE order_no = ?", data.AppOrderNo)
				if e != nil {
					g.Log().Info("RunSubLaunchpadPayTask err:%v", e)
					return nil
				}
			}
			e = tx.Commit()
		} else {
			g.Log().Info("RunSubLaunchpadPayTask json.Unmarshal err:%v", err)
		}
		return nil
	})
}
