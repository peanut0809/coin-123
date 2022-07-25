package service

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/model"
	"time"
)

type seckillActivity struct {
}

var SeckillActivity = new(seckillActivity)

func (s *seckillActivity) GetValidDetail(alias string) (ret model.SeckillActivityFull, err error) {
	var as *model.SeckillActivity
	now := time.Now()
	err = g.DB().Model("seckill_activity").Where("alias = ? AND start_time < ?", alias, now).Scan(&as)
	if err != nil {
		return
	}
	if as == nil {
		err = fmt.Errorf("活动不存在")
		return
	}
	ret.SeckillActivity = as
	if now.Unix() > as.ActivityStartTime.Unix() && now.Unix() < as.ActivityEndTime.Unix() {
		ret.Status = model.SeckillActivityStatus_Ing
	} else {
		if now.Unix() < as.ActivityStartTime.Unix() {
			ret.Status = model.SeckillActivityStatus_Wait_Start
			ret.LastSec = as.ActivityStartTime.Unix() - now.Unix()
		} else {
			ret.Status = model.SeckillActivityStatus_End
		}
	}
	//provider.Asset.
	return
}
