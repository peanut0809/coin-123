package task

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/client"
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"meta_launchpad/provider"
	"meta_launchpad/service"
	"time"
)

func RunSeckillOrderPayTask() {
	queue, err := client.GetQueue(client.QueueConfig{
		QueueName:  "payment.success.launchpad_seckill",
		Exchange:   "payment.success.launchpad_seckill",
		RoutingKey: "",
		Kind:       "fanout",
		MqUrl:      g.Cfg().GetString("rabbitmq.default.link"),
	})
	if err != nil {
		g.Log().Error("启动mq失败：" + err.Error())
		return
	}
	queue.RunConsume(func(msg []byte) error {
		g.Log().Info("RunSeckillOrderPayTask receive:%s", gconv.String(msg))
		var data RunSubLaunchpadPayData
		err = json.Unmarshal(msg, &data)
		if err == nil {
			orderInfo, err := service.SeckillOrder.GetByOrderNos([]string{data.AppOrderNo})
			if err != nil {
				g.Log().Errorf("RunSeckillOrderPayTask err:%v", err)
				return nil
			}
			if len(orderInfo) == 0 {
				g.Log().Errorf("RunSeckillOrderPayTask err:%v", err)
				return nil
			}
			activityInfo, err := service.SeckillActivity.GetSimpleDetail(orderInfo[0].Aid)
			if err != nil {
				g.Log().Errorf("RunSubLaunchpadPayTask err:%v", err)
				return nil
			}
			tx, e := g.DB().Begin()
			if e != nil {
				g.Log().Info("RunSeckillOrderPayTask err:%v", e)
			} else {
				r, e := tx.Exec("UPDATE seckill_orders SET status = 2,pay_at = ?,pay_method = ? WHERE order_no = ?", gtime.Now(), data.PayType, data.AppOrderNo)
				if e != nil {
					g.Log().Info("RunSeckillOrderPayTask err:%v", e)
					return nil
				}
				affectedNum, _ := r.RowsAffected()
				if affectedNum != 1 {
					tx.Rollback()
					return nil
				}
				_, e = tx.Exec("DELETE FROM seckill_wait_pay_orders WHERE order_no = ?", data.AppOrderNo)
				if e != nil {
					g.Log().Info("RunSeckillOrderPayTask err:%v", e)
					tx.Rollback()
					return nil
				}
				e = tx.Commit()
				if e != nil {
					g.Log().Info("RunSeckillOrderPayTask err:%v", e)
					tx.Rollback()
					return nil
				}
				//发放资产
				publishSuccess := false
				for i := 0; i < 3; i++ {
					err = provider.Asset.PublishAssetWithTemplateId(&map[string]interface{}{
						"appId":      activityInfo.AppId,
						"templateId": activityInfo.TemplateId,
						"num":        orderInfo[0].Num,
						"userId":     orderInfo[0].UserId,
						"optType":    "LAUNCHPAD",
						"optRemark":  "元初发射台秒杀发放资产",
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
					_, e = g.DB().Exec("UPDATE seckill_orders SET publish_status = 1 WHERE order_no = ?", data.AppOrderNo)
					if e != nil {
						g.Log().Info("RunSeckillOrderPayTask err:%v", e)
						return nil
					}
					//通知钱包
					if data.PayType == "wallet_pay" {
						service.NoticeWallet(service.NoticeWalletReq{
							FromUserId: orderInfo[0].UserId,
							ToUserId:   "B",
							TotalFee:   orderInfo[0].Price,
						})
					}
				} else {
					_, e = g.DB().Exec("UPDATE seckill_orders SET publish_status = 2 WHERE order_no = ?", data.AppOrderNo)
					if e != nil {
						g.Log().Info("RunSeckillOrderPayTask err:%v", e)
						return nil
					}
				}
			}
		} else {
			g.Log().Info("RunSeckillOrderPayTask json.Unmarshal err:%v", err)
		}
		return nil
	})
}
