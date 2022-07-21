package provider

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	"context"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

type payment struct {
}

var Payment = new(payment)

type CreateOrderReq struct {
	Subject            string      `json:"subject"`
	Description        string      `json:"description"`
	PayAmount          int         `json:"pay_amount"`
	UserId             string      `json:"user_id"`
	AppType            string      `json:"app_type"`
	SuccessRedirectUrl string      `json:"success_redirect_url"`
	ExitRedirectUrl    string      `json:"exit_redirect_url"`
	ClientIp           string      `json:"client_ip"`
	Extra              string      `json:"extra"`
	PayExpire          *gtime.Time `json:"pay_expire"`
	AppOrderNo         string      `json:"app_order_no"`
	PublisherId        string      `json:"publisherId"`
	PlatformAppId      string      `json:"platformAppId"`
}

func (c *payment) CreateOrder(req *CreateOrderReq) (err error) {
	_, err = utils.SendJsonRpc(context.Background(), "payment", "PayOrderBeforehand.CreatePayOrderBeforehand", req)
	if err != nil {
		g.Log().Errorf("下单失败：%v,data:%+v", err, req)
		return
	}
	return
}
