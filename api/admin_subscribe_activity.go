package api

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/api"
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	"github.com/gogf/gf/net/ghttp"
	"github.com/shopspring/decimal"
	"meta_launchpad/model"
	"meta_launchpad/service"
)

type adminSubscribeActivity struct {
	api.CommonBase
}

var AdminSubscribeActivity = new(adminSubscribeActivity)

func (s *adminSubscribeActivity) Create(r *ghttp.Request) {
	var req model.CreateSubscribeActivityReq
	err := r.Parse(&req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	if req.Name == "" || req.CoverImgUrl == "" || req.AppId == "" || req.TemplateId == "" || req.ActivityIntro == "" || req.PriceYuan == "" || req.SumNum < 0 {
		s.FailJsonExit(r, "缺少必要参数")
		return
	}
	if req.ActivityStartTime == nil || req.ActivityEndTime == nil || req.OpenAwardTime == nil || req.PayEndTime == nil {
		s.FailJsonExit(r, "缺少必要参数")
		return
	}
	if req.ActivityStartTime.Unix() > req.ActivityEndTime.Unix() {
		s.FailJsonExit(r, "开始时间不能大于结束时间")
		return
	}
	if req.ActivityEndTime.Unix() > req.OpenAwardTime.Unix() {
		s.FailJsonExit(r, "结束时间不能大于开奖时间")
		return
	}
	if req.OpenAwardTime.Unix() > req.PayEndTime.Unix() {
		s.FailJsonExit(r, "开奖时间不能大于支付截止时间")
		return
	}
	decimalValue, err := decimal.NewFromString(req.PriceYuan)
	if err != nil {
		s.FailJsonExit(r, "价格参数错误")
		return
	}
	if req.ActivityType != 1 && req.ActivityType != 2 {
		s.FailJsonExit(r, "活动类型错误")
		return
	}
	decimalValue = decimalValue.Mul(decimal.NewFromFloat(100))
	priceInt := decimalValue.IntPart()
	if priceInt < 0 {
		s.FailJsonExit(r, "价格参数错误")
		return
	}
	req.Price = int(priceInt)
	req.Alias = utils.RandString(6)

	err = service.AdminSubscribeActivity.Create(req.SubscribeActivity, req.Condition)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r)
}
