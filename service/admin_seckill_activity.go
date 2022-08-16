package service

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/model"
)

type adminSecKillActivity struct {
}

var AdminSecKillActivity = new(adminSecKillActivity)

func (s *adminSecKillActivity) Create(in model.SeckillActivity) (err error) {
	_, err = g.DB().Model("seckill_activity").Insert(&in)
	if err != nil {
		return
	}
	return
}

func (s *adminSecKillActivity) Update(in model.SeckillActivity) (err error) {
	updateMap := g.Map{
		"name":                in.Name,
		"app_id":              in.AppId,
		"template_id":         in.TemplateId,
		"price":               in.Price,
		"sum_num":             in.SumNum,
		"remain_num":          in.SumNum,
		"cover_img_url":       in.CoverImgUrl,
		"limit_buy":           in.LimitBuy,
		"activity_intro":      in.ActivityIntro,
		"activity_start_time": in.ActivityStartTime,
		"activity_end_time":   in.ActivityEndTime,
	}
	_, err = g.DB().Model("seckill_activity").Data(updateMap).Where("id = ? AND publisher_id = ?", in.Id, in.PublisherId).Update()
	if err != nil {
		return
	}
	return
}

func (s *adminSecKillActivity) Detail(id int, publisherId string) (ret model.CreateSeckillActivityReq, err error) {
	var activityInfo *model.SeckillActivity
	err = g.DB().Model("seckill_activity").Where("id = ? AND publisher_id = ?", id, publisherId).Scan(&activityInfo)
	if err != nil {
		return
	}
	if activityInfo == nil {
		err = fmt.Errorf("活动不存在")
		return
	}
	ret.SeckillActivity = *activityInfo
	ret.PriceYuan = fmt.Sprintf("%.2f", float64(activityInfo.Price)/100)
	return
}
