package provider

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	"context"
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
	DetailImg string `json:"detail_img"`
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
		AppId       string `json:"appId"`
		CnName      string `json:"cnName"`
		PublisherId string `json:"publisherId"`
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

type MwdAppsPublisher struct {
	Id          int    `orm:"id" json:"id"`
	Name        string `orm:"name" json:"name"`
	Icon        string `orm:"icon" json:"icon"`
	PublisherId string `orm:"publisher_id" json:"publisherId"`
	Desc        string `orm:"desc" json:"desc"`
	UserId      int    `orm:"user_id" json:"userId"`
	IsDeleted   int    `orm:"is_deleted" json:"isDeleted"`
	WebsiteUrl  string `json:"websiteUrl"`
	ChainName   string `json:"chainName"`
	ChainAddr   string `json:"chainAddr"`
	ChainType   int    `json:"chainType"`
}

func (s *developer) GetPublishInfo(publisherId string) (ret MwdAppsPublisher, err error) {
	params := &map[string]interface{}{
		"publisherId": publisherId,
	}
	err = utils.SendJsonRpcScan(context.Background(), "developer", "Publisher.GetPublisherById", params, &ret)
	if err != nil {
		return
	}
	return
}

func (s *developer) GetPublisherByIds(publisherId []string) (ret map[string]MwdAppsPublisher, err error) {
	params := &map[string]interface{}{
		"publisherIds": publisherId,
	}
	err = utils.SendJsonRpcScan(context.Background(), "developer", "Publisher.GetPublisherByIds", params, &ret)
	if err != nil {
		return
	}
	return
}
