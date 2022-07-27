package task

import (
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/service"
)

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
		e = service.SeckillActivity.UpdateRemain(tx, v.Aid, v.Num)
		if e != nil {
			tx.Rollback()
			g.Log().Errorf("CheckSeckillOrderTimeout err:%v", e)
			return
		}
		e = service.SeckillUserBnum.UpdateRemain(tx, v.UserId, v.Aid, v.Num)
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
	}
}
