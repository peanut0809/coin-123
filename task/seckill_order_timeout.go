package task

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/cache"
	"meta_launchpad/service"
	"time"
)

const TASK_CheckSeckillOrderTimeoutTask = "CheckSeckillOrderTimeoutTask"

//检查超时未支付
func CheckSeckillOrderTimeoutTask() {
	cache.DistributedUnLock(TASK_CheckSeckillOrderTimeoutTask)
	for {
		lock := cache.DistributedLock(TASK_CheckSeckillOrderTimeoutTask)
		if lock {
			CheckSeckillOrderTimeout()
			cache.DistributedUnLock(TASK_CheckSeckillOrderTimeoutTask)
		}
		time.Sleep(time.Second * 10)
	}
}

func CheckSeckillOrderTimeout() {
	worders, err := service.SeckillWaitPayOrder.GetWaitPayOrder()
	if err != nil {
		g.Log().Errorf("CheckSeckillOrderTimeout err:%v", err)
		return
	}
	if len(worders) == 0 {
		return
	}
	orderNos := make([]string, 0)
	for _, v := range worders {
		orderNos = append(orderNos, v.OrderNo)
	}
	orders, err := service.SeckillOrder.GetByOrderNos(orderNos)
	if err != nil {
		g.Log().Errorf("CheckSeckillOrderTimeout err:%v", err)
		return
	}
	err = service.SeckillOrder.UpdateOrderNosStatus(orderNos, 3)
	if err != nil {
		g.Log().Errorf("CheckSeckillOrderTimeout err:%v", err)
		return
	}
	for _, v := range orders {
		tx, e := g.DB().Begin()
		if e != nil {
			g.Log().Errorf("CheckSeckillOrderTimeout err:%v", e)
			return
		}
		e = service.SeckillOrder.CancelHandel(tx, v.Aid, v.Num, v.OrderNo, v.UserId)
		if e != nil {
			tx.Rollback()
			g.Log().Errorf("CheckSeckillOrderTimeout err:%v", e)
			return
		}
		e = tx.Commit()
		if e != nil {
			tx.Rollback()
			g.Log().Errorf("CheckSeckillOrderTimeout err:%v", e)
			return
		}
		//设置处罚时间
		g.Redis().Do("SET", fmt.Sprintf(cache.SECKILL_DISCIPLINE, v.UserId), 1, "ex", 3600*24*30)
	}
}
