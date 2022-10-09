package api

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/api"
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	"bytes"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/parnurzeal/gorequest"
	"github.com/xuri/excelize/v2"
	"meta_launchpad/model"
	"meta_launchpad/provider"
	"meta_launchpad/service"
	"strings"
)

type drop struct {
	api.CommonBase
}

var Drop = new(drop)

func (s *drop) GetDetailRecordList(r *ghttp.Request) {
	pageNum := r.GetInt("pageNum", 1)
	pageSize := r.GetInt("pageSize", 20)
	dropId := r.GetInt("dropId")
	ret, err := service.Drop.GetDetailRecordList(pageNum, pageSize, dropId)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *drop) GetRecordList(r *ghttp.Request) {
	pageNum := r.GetInt("pageNum", 1)
	pageSize := r.GetInt("pageSize", 20)
	createStartTime := r.GetString("createStartTime")
	createEndTime := r.GetString("createEndTime")
	searchVal := r.GetString("searchVal")
	ret, err := service.Drop.GetRecordList(pageNum, pageSize, createStartTime, createEndTime, searchVal)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *drop) Create(r *ghttp.Request) {
	var req model.DropRecordReq
	err := r.Parse(&req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	if req.Name == "" || req.AppId == "" || req.TemplateId == "" || req.Num <= 0 || req.NfrSec < 0 {
		s.FailJsonExit(r, "参数错误")
		return
	}
	appInfo, e := provider.Developer.GetAppInfo(req.AppId)
	if err != nil {
		s.FailJsonExit(r, e.Error())
		return
	}
	if appInfo.Data.PublisherId != s.GetPublisherId(r) {
		s.FailJsonExit(r, "无权操作")
		return
	}
	if len(req.PhoneArr) == 0 {
		if req.ExcelFile == "" {
			s.FailJsonExit(r, "参数错误")
			return
		}
		_, bs, errs := gorequest.New().Get(req.ExcelFile).EndBytes()
		if len(errs) != 0 {
			s.FailJsonExit(r, "文件地址错误")
			return
		}
		f, e := excelize.OpenReader(bytes.NewReader(bs))
		if e != nil {
			s.FailJsonExit(r, "文件读取错误")
			return
		}
		rows, e := f.GetRows("Sheet1")
		if e != nil {
			s.FailJsonExit(r, "文件读取错误")
			return
		}
		for _, v := range rows {
			req.PhoneArr = append(req.PhoneArr, strings.TrimSpace(v[0]))
		}
	}
	if len(req.PhoneArr) <= 0 {
		s.FailJsonExit(r, "缺少手机号")
		return
	}
	phoneMap := make(map[string]int)
	for _, v := range req.PhoneArr {
		phoneMap[v] = 1
	}
	if len(phoneMap) != len(req.PhoneArr) {
		s.FailJsonExit(r, "手机号重复")
		return
	}
	if len(req.PhoneArr) > 5000 {
		s.FailJsonExit(r, "手机号不能超过5000个")
		return
	}
	if req.Num > 10 {
		s.FailJsonExit(r, "空投资产单次不能超过10个")
		return
	}
	req.Phones = strings.Join(req.PhoneArr, ",")
	req.OrderNo = utils.Generate()
	locKey := fmt.Sprintf("meta_launchpad:drop:asset")
	re, e := g.Redis().Do("SET", locKey, 1, "ex", 60, "nx")
	if fmt.Sprintf("%v", re) == "OK" && e == nil {
		defer g.Redis().Do("DEL", locKey)
		err = service.Drop.Create(req.DropRecord)
		if err != nil {
			s.FailJsonExit(r, err.Error())
			return
		}
	} else {
		s.FailJsonExit(r, "操作太快了，请稍后重试")
		return
	}
	s.SusJsonExit(r)
}
