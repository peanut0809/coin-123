package task

import (
	"context"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/cache"
	"meta_launchpad/model"
	"meta_launchpad/service"
	"time"
)

const TASKCheckEquityOrderTimeoutTask = "CheckEquityOrderTimeoutTask"

// 检查超时未支付
func CheckEquityOrderTimeoutTask() {
	cache.DistributedUnLock(TASKCheckEquityOrderTimeoutTask)
	for {
		lock := cache.DistributedLock(TASKCheckEquityOrderTimeoutTask)
		if lock {
			CheckEquityOrderTimeout()
			cache.DistributedUnLock(TASKCheckEquityOrderTimeoutTask)
		}
		time.Sleep(time.Second * 10)
	}
}

func CheckEquityOrderTimeout() {
	workers, err := service.EquityOrder.GetWaitPayOrder()
	if err != nil {
		g.Log().Errorf("CheckEquityOrderTimeout err:%v", err)
		return
	}
	if len(workers) == 0 {
		return
	}
	for _, v := range workers {
		orderInfo, err := service.EquityOrder.GetInfoByOrderNo(v.OrderNo)
		if err != nil {
			g.Log().Errorf("CheckEquityOrderTimeout err:%v", err)
			return
		}
		err = g.DB().Transaction(context.Background(), func(ctx context.Context, tx *gdb.TX) error {
			err = service.EquityOrder.UpdateOrderNoStatus(v.OrderNo, model.TIMEOUT)
			if err != nil {
				g.Log().Errorf("CheckEquityOrderTimeout err:%v", err)
				return err
			}
			return service.EquityOrder.InventoryRollback(tx, orderInfo.ActivityId, orderInfo.Num)
		})
		if err != nil {
			g.Log().Errorf("CheckEquityOrderTimeout err:%v", err)
			return
		}
	}
}
