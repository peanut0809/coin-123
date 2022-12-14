package task

import (
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/cache"
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
		err := service.EquityOrder.Cancel(v.UserId, v.OrderNo)
		if err != nil {
			g.Log().Errorf("CheckEquityOrderTimeout err:%v", err)
			return
		}
	}
}
