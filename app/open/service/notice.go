package service

import (
	"encoding/json"
	"time"

	"brq5j1d.gfanx.pro/meta_cloud/meta_service/app/third/model"
	"brq5j1d.gfanx.pro/meta_cloud/meta_service/app/third/utils"

	"github.com/gogf/gf/frame/g"
	"github.com/parnurzeal/gorequest"
)

type noticeService struct {
}

var NoticeService = new(noticeService)

//func init() {
//	NoticeService.RetryNoticeThird("202201112132030941540002")
//}

func (s *noticeService) NoticeThird(orderInfo model.ThirdOrder) bool {
	req, err := utils.MakeSign(orderInfo.Appid, orderInfo)
	if err != nil {
		g.Log("pay").Errorf("NoticeThird MakeSign err:%v", err)
		return false
	}
	reqByte, err := json.Marshal(req)
	if err != nil {
		g.Log("pay").Errorf("NoticeThird json.Marshal err:%v", err)
		return false
	}
	g.Log("pay").Infof("请求third服务，url:%v\nbody=> %v", orderInfo.NotifyUrl, string(reqByte))
	_, resp, errs := gorequest.New().Post(orderInfo.NotifyUrl).Timeout(time.Second*10).AppendHeader("Content-Type", "application/json; encoding=utf-8").SendString(string(reqByte)).EndBytes()
	g.Log("pay").Infof("third返回结果 => resp:%v", string(resp))
	if len(errs) != 0 || "SUCCESS" != string(resp) {
		return false
	}
	return true
}

func (s *noticeService) RetryNoticeThird(orderNo string) (err error) {
	ins := make([]model.ThirdRetryNotice, 0)
	for i := 2; i < 6; i++ {
		ins = append(ins, model.ThirdRetryNotice{
			OrderNo:  orderNo,
			NoticeAt: time.Now().Add(time.Minute * time.Duration((i-1)*5)),
			Num:      i,
		})
	}
	_, err = g.DB("meta_world").Model("third_retry_notices").Insert(&ins)
	if err != nil {
		return
	}
	return
}

func (s *noticeService) GetWaitNotices() (ret []model.ThirdRetryNotice) {
	_ = g.DB("meta_world").Model("third_retry_notices").Limit(100).Scan(&ret)
	return
}

func (s *noticeService) DelWaitNotices(orderNo string) (err error) {
	_, err = g.DB("meta_world").Model("third_retry_notices").Where("order_no = ?", orderNo).Delete()
	return
}

func (s *noticeService) DelWaitNoticesById(id int) (err error) {
	_, err = g.DB("meta_world").Model("third_retry_notices").Where("id = ?", id).Delete()
	return
}
