package service

import (
	"encoding/json"
	"fmt"
	"time"

	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"

	"brq5j1d.gfanx.pro/meta_cloud/meta_service/app/third/model"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
)

type orderService struct {
}

var OrderService = new(orderService)

func (s *orderService) Create(in model.ThirdOrder) (out model.ThirdOrder, err error) {
	_, err = g.DB("meta_world").Model("third_orders").Insert(&in)
	out = in
	return
}

func (s *orderService) MakeOrderSign(appid, orderNo, userId string) string {
	src := fmt.Sprintf(`{"appid":"%s","orderNo":"%s","userId":"%s"}`, appid, orderNo, userId)
	sign := utils.AesEncrypt(src, g.Cfg().GetString("gToken.ucenter.EncryptKey"))
	g.Log().Infof("生成订单签名：srs：%s --- sing：%s", src, sign)
	return sign
}

func (s *orderService) DecodeOrderSign(orderSign string) (ret map[string]string, err error) {
	//src := fmt.Sprintf(`{"appid":"%s","orderNo":"%s"}`, appid, orderNo)
	orderSignByte, err := utils.Base64URLDecode(orderSign)
	if err != nil {
		return
	}
	orgData, err := utils.AesDecrypt(orderSignByte, []byte(g.Cfg().GetString("gToken.ucenter.EncryptKey")))
	if err != nil {
		return
	}
	err = json.Unmarshal(orgData, &ret)
	return
}

//更新订单信息
func (s *orderService) UpdateOrderInfo(orderNo string, notifyStatus, notifyNum int) (err error) {
	if notifyStatus == 1 {
		_, err = g.DB("meta_world").Model("third_orders").Data(g.Map{
			"notify_status": notifyStatus,
			"notify_num":    notifyNum,
			"notify_time":   time.Now(),
		}).Where("order_no = ?", orderNo).Update()
		return
	}
	if notifyStatus == 2 {
		_, err = g.DB("meta_world").Model("third_orders").Data(g.Map{
			"notify_status": notifyStatus,
		}).Where("order_no = ?", orderNo).Update()
		return
	}
	if notifyStatus == 3 {
		_, err = g.DB("meta_world").Model("third_orders").Data(g.Map{
			"notify_status": notifyStatus,
			"notify_num":    notifyNum,
		}).Where("order_no = ?", orderNo).Update()
		return
	}
	return
}

//更新订单为已支付状态
func (s *orderService) Paid(tx *gdb.TX, orderNo string) (err error) {
	_, err = tx.Exec("UPDATE third_orders SET status = 1,pay_time = ? WHERE order_no = ?", time.Now(), orderNo)
	if err != nil {
		return
	}
	return
}
