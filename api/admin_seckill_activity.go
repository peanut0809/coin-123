package api

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/api"
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	"github.com/gogf/gf/net/ghttp"
	"github.com/shopspring/decimal"
	"meta_launchpad/model"
	"meta_launchpad/service"
)

type adminSeckillActivity struct {
	api.CommonBase
}

var AdminSeckillActivity = new(adminSeckillActivity)

func (s *adminSeckillActivity) Create(r *ghttp.Request) {
	var req model.CreateSeckillActivityReq
	err := r.Parse(&req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	if req.Name == "" || req.AppId == "" || req.TemplateId == "" || req.SumNum <= 0 || req.CoverImgUrl == "" || req.LimitBuy <= 0 || req.ActivityIntro == "" || req.ActivityStartTime == nil || req.ActivityEndTime == nil {
		s.FailJsonExit(r, "参数错误")
		return
	}
	decimalValue, err := decimal.NewFromString(req.PriceYuan)
	if err != nil {
		s.FailJsonExit(r, "价格参数错误")
		return
	}
	decimalValue = decimalValue.Mul(decimal.NewFromFloat(100))
	priceInt := decimalValue.IntPart()
	if priceInt <= 0 {
		s.FailJsonExit(r, "价格参数错误")
		return
	}
	if req.ActivityStartTime.Unix() > req.ActivityEndTime.Unix() {
		s.FailJsonExit(r, "活动开始时间不能大于结束时间")
		return
	}
	req.Alias = utils.RandString(6)
	req.PublisherId = s.GetPublisherId(r)
	req.Price = int(priceInt)
	req.RemainNum = req.SumNum
	if req.Id == 0 {
		err = service.AdminSecKillActivity.Create(req.SeckillActivity)
		if err != nil {
			s.FailJsonExit(r, err.Error())
			return
		}
	} else {
		//service.SeckillOrder.Create()
		err = service.AdminSecKillActivity.Update(req.SeckillActivity)
		if err != nil {
			s.FailJsonExit(r, err.Error())
			return
		}
	}
	s.SusJsonExit(r)
}

func (s *adminSeckillActivity) Detail(r *ghttp.Request) {
	id := r.GetQueryInt("id")
	ret, err := service.AdminSecKillActivity.Detail(id, s.GetPublisherId(r))
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}
