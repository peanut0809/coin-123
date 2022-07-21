package task

import (
	"encoding/json"
	"fmt"
	"meta_open_sdk/app/open/service"

	"brq5j1d.gfanx.pro/meta_cloud/meta_service/app/third/model"

	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/client"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
)

type PaymentQueueData struct {
	AppOrderNo string `json:"app_order_no"`
}

func PaymentQueue() {
	queue, err := client.GetQueue(client.QueueConfig{
		QueueName:  "payment.success.sdk",
		Exchange:   "payment.success.sdk",
		RoutingKey: "",
		Kind:       "fanout",
		MqUrl:      g.Cfg().GetString("rabbitmq.default.link"),
	})
	if err != nil {
		g.Log().Error(err)
		return
	}
	queue.RunConsume(func(msg []byte) error {
		g.Log().Info("PaymentQueue receive:%s", gconv.String(msg))
		var data PaymentQueueData
		json.Unmarshal(msg, &data)

		//查出订单信息
		// var orderInfo model.ThirdOrderFull
		orderInfo, err := service.PayService.GetOrderDetailByOrderNo(data.AppOrderNo)
		if err != nil {
			return err
		}

		err = UpdateOrder(orderInfo)
		if err != nil {
			return err
		}
		err = NotifyThird(orderInfo)
		if err != nil {
			return err
		}

		return nil
	})
}

func UpdateOrder(orderInfo model.ThirdOrderFull) (err error) {
	//防止同一个订单并发
	re, e := g.Redis().Do("SET", fmt.Sprintf("order_%s", orderInfo.OrderNo), 1, "ex", 60, "nx")
	if fmt.Sprintf("%v", re) == "OK" && e == nil {
		//释放锁
		defer g.Redis().Do("DEL", fmt.Sprintf("order_%s", orderInfo.OrderNo))
		tx, e := g.DB("meta_world").Begin()
		if e != nil {
			err = e
			g.Log("pay").Errorf("meta_world begin:%v", err)
			return
		}
		//更新订单状态
		err = service.OrderService.Paid(tx, orderInfo.OrderNo)
		if err != nil {
			g.Log("pay").Errorf("meta_world UPDATE third_orders err:%v", err)
			tx.Rollback()
			return
		}
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			g.Log("pay").Errorf("tx.Commit err:%v", err)
			return
		}
	}
	err = fmt.Errorf("订单已被支付过")
	return
}

func NotifyThird(orderInfo model.ThirdOrderFull) (err error) {

	//通知第三方服务
	go func() {
		isSuccess := service.NoticeService.NoticeThird(orderInfo.ThirdOrder)
		if isSuccess { //如果一次通知成功，更新订单信息
			err = service.OrderService.UpdateOrderInfo(orderInfo.OrderNo, 1, 1)
			if err != nil {
				g.Log("pay").Errorf("thirdService.OrderService.UpdateOrderInfo err:%v", err)
				return
			}
			return
		}
		//将订单加入重试队列
		err = service.OrderService.UpdateOrderInfo(orderInfo.OrderNo, 2, 0)
		if err != nil {
			g.Log("pay").Errorf("thirdService.OrderService.UpdateOrderInfo err:%v", err)
			return
		}
		err = service.NoticeService.RetryNoticeThird(orderInfo.ThirdOrder.OrderNo)
		if err != nil {
			g.Log("pay").Errorf("thirdService.NoticeService.RetryNoticeThird err:%v", err)
			return
		}
	}()

	return
}
