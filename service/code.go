package service

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/parnurzeal/gorequest"
	"time"
)

type sms struct {
}

var Sms = new(sms)

type sendCodeReq struct {
	AccessKey            string            `json:"accessKey"`
	AccessSecret         string            `json:"accessSecret"`
	ClassificationSecret string            `json:"classificationSecret"`
	SignCode             string            `json:"signCode"`
	TemplateCode         string            `json:"templateCode"`
	Phone                string            `json:"phone"`
	Params               map[string]string `json:"params"`
}

type sendCodeRet struct {
	BusinessData struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	} `json:"BusinessData"`
}

func (s *sms) SendSms(phone string, signCode string, templateCode string, classificationSecret string, params map[string]string) (err error) {
	req := sendCodeReq{
		AccessKey:            g.Cfg().GetString("sms.accessKey"),
		AccessSecret:         g.Cfg().GetString("sms.accessSecret"),
		ClassificationSecret: classificationSecret,
		SignCode:             signCode,
		TemplateCode:         templateCode,
		Phone:                phone,
		Params:               params,
	}
	ret := sendCodeRet{}
	sendbByte, _ := json.Marshal(req)
	_, _, errs := gorequest.New().Post(g.Cfg().GetString("sms.url")).Timeout(time.Second*10).AppendHeader("Content-Type", "application/json; encoding=utf-8").SendString(string(sendbByte)).EndStruct(&ret)
	if len(errs) != 0 {
		err = fmt.Errorf("失败-1")
		g.Log().Errorf("发送验证码失败：%v", errs[0])
		return
	}
	if ret.BusinessData.Code != 10000 {
		g.Log().Errorf("发送验证码失败：%s", ret.BusinessData.Msg)
		err = fmt.Errorf("获取验证码次数超限，0点后可重新获取")
		return
	}
	return
}
