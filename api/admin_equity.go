package api

import (
	"fmt"
	"meta_launchpad/model"
	"meta_launchpad/service"

	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/api"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gvalid"
	"github.com/shopspring/decimal"
)

type adminEquity struct {
	api.CommonBase
}

var AdminEquity = new(adminEquity)

func (s *adminEquity) Create(r *ghttp.Request) {
	var req model.CreateEquityActivityReq
	err := r.Parse(&req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}

	if req.Name == "" || req.TimeType <= 0 || req.SubLimitType <= 0 || req.TemplateId == "" || req.Number <= 0 || req.LimitBuy <= 0 || req.LimitType <= 0 || req.ActivityStartTime == nil || req.ActivityEndTime == nil {
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
	req.Price = int(priceInt)

	err = service.AdminEquity.Create(req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r)
}

// 用户导入解析
func (s *adminEquity) Import(r *ghttp.Request) {
	var req model.CreateEquityActivityReq
	result, err := service.AdminEquity.HandelExcelUser(req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, result)
}

// 获取详情
func (s *adminEquity) Item(r *ghttp.Request) {
	templateId := r.GetQueryString("templateId")
	if templateId == "" {
		s.FailJsonExit(r, "活动标识为空")
		return
	}
	ret, err := service.AdminEquity.Item(templateId)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

//Invalid 失效下架活动
// 获取详情
func (s *adminEquity) Invalid(r *ghttp.Request) {

}

// 获取专属活动用户明细
func (s *adminEquity) UserItems(r *ghttp.Request) {

	var req model.EquityUserReq
	if err := r.Parse(&req); err != nil {
		s.FailJsonExit(r, err.(gvalid.Error).FirstString())
		return
	}
	//获取参数
	fmt.Println(req)
	if req.EquityId <= 0 {
		s.FailJsonExit(r, "活动标识为空")
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	ret, err := service.AdminEquity.UserItems(req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}
