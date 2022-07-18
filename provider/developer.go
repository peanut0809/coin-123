package provider

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/parnurzeal/gorequest"
)

type developer struct {
}

var Developer = new(developer)

type AssetsTemplateRet struct {
	Code int          `json:"code"`
	Data TemplateInfo `json:"data"`
	Msg  string       `json:"msg"`
}

type TemplateInfo struct {
	CateList []struct {
		CnName string `json:"cnName"`
	} `json:"cate_list"`
}

func (s *developer) GetAssetsTemplate(appId string, templateId string) (ret TemplateInfo, err error) {
	developerHost := g.Cfg().GetString("developer.host")
	var info AssetsTemplateRet
	_, _, errs := gorequest.New().Post(developerHost + "/out/tpl-detail").SendString(fmt.Sprintf(`{"appId":"%s","tplId":"%s"}`, appId, templateId)).EndStruct(&info)
	if len(errs) != 0 {
		err = errs[0]
		return
	}
	if info.Code != 200 {
		err = fmt.Errorf(info.Msg)
		return
	}
	ret = info.Data
	return
}

type AppServerDetail struct {
	Code int `json:"code"`
	Data struct {
		AppId  string `json:"appId"`
		CnName string `json:"cnName"`
	} `json:"data"`
}

func (s *developer) GetAppInfo(appId string) (ret AppServerDetail, err error) {
	developHost := g.Cfg().GetString("developer.host")
	rurl := fmt.Sprintf("%s/out/app-server-detail", developHost)
	_, _, errs := gorequest.New().Post(rurl).AppendHeader("Content-Type", "application/json; encoding=utf-8").SendString(fmt.Sprintf(`{"appId":"%s"}`, appId)).EndStruct(&ret)
	if len(errs) != 0 {
		err = errs[0]
		return
	}
	return
}
