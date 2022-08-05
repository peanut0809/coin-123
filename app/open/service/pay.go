package service

import (
	"fmt"

	"brq5j1d.gfanx.pro/meta_cloud/meta_service/app/third/model"

	"brq5j1d.gfanx.pro/meta_cloud/meta_service/app/assets/provider"

	"github.com/gogf/gf/frame/g"
)

var PayService = new(payService)

type payService struct {
}

//获取订单信息
func (s *payService) GetOrderDetail(appid, userId, orderSign string, cny int) (ret model.ThirdOrderFull, err error) {
	orgData, err := OrderService.DecodeOrderSign(orderSign)
	if err != nil {
		g.Log().Errorf("订单解密失败:%s err:%v", orderSign, err)
		err = fmt.Errorf("订单解密失败")
		return
	}
	if appid != orgData["appid"] {
		err = fmt.Errorf("无权查看订单信息")
		return
	}
	if orgData["userId"] != userId {
		err = fmt.Errorf("无权查看订单信息")
		return
	}
	err = g.DB("meta_world").Model("third_orders").Where("order_no = ?", orgData["orderNo"]).Scan(&ret.ThirdOrder)
	if err != nil {
		return
	}
	ret.SumCny = fmt.Sprintf("%.2f", float64(cny)/100)
	ret.CanPay = (cny >= ret.ThirdOrder.Amount)
	ret.AmountCny = fmt.Sprintf("%.2f", float64(ret.ThirdOrder.Amount)/100)
	appInfo, _ := provider.DeveloperProvider.GetAppInfo(appid)
	ret.GameName = appInfo.CnName
	return
}

//获取订单信息
func (s *payService) GetOrderDetailByOrderNo(orderSign string) (ret model.ThirdOrderFull, err error) {
	orgData, err := OrderService.DecodeOrderSign(orderSign)
	if err != nil {
		g.Log().Errorf("订单解密失败:%s err:%v", orderSign, err)
		err = fmt.Errorf("订单解密失败")
		return
	}
	// if appid != orgData["appid"] {
	// 	err = fmt.Errorf("无权查看订单信息")
	// 	return
	// }
	// if orgData["userId"] != userId {
	// 	err = fmt.Errorf("无权查看订单信息")
	// 	return
	// }
	err = g.DB("meta_world").Model("third_orders").Where("order_no = ?", orgData["orderNo"]).Scan(&ret.ThirdOrder)
	if err != nil {
		return
	}
	// ret.SumCny = fmt.Sprintf("%.2f", float64(cny)/100)
	// ret.CanPay = (cny >= ret.ThirdOrder.Amount)
	ret.AmountCny = fmt.Sprintf("%.2f", float64(ret.ThirdOrder.Amount)/100)
	// appInfo, _ := provider.DeveloperProvider.GetAppInfo(appid)
	// ret.GameName = appInfo.CnName
	return
}

//获取订单信息
func (s *payService) GetOrderDetailByOrderNoNoSign(orderSign string) (ret model.ThirdOrderFull, err error) {
	// if appid != orgData["appid"] {
	// 	err = fmt.Errorf("无权查看订单信息")
	// 	return
	// }
	// if orgData["userId"] != userId {
	// 	err = fmt.Errorf("无权查看订单信息")
	// 	return
	// }
	err = g.DB("meta_world").Model("third_orders").Where("order_no = ?", orderSign).Scan(&ret.ThirdOrder)
	if err != nil {
		return
	}
	// ret.SumCny = fmt.Sprintf("%.2f", float64(cny)/100)
	// ret.CanPay = (cny >= ret.ThirdOrder.Amount)
	ret.AmountCny = fmt.Sprintf("%.2f", float64(ret.ThirdOrder.Amount)/100)
	// appInfo, _ := provider.DeveloperProvider.GetAppInfo(appid)
	// ret.GameName = appInfo.CnName
	return
}
