package service

import (
	"bytes"
	"fmt"
	"meta_launchpad/model"
	"meta_launchpad/provider"
	"strconv"
	"strings"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/parnurzeal/gorequest"
	"github.com/xuri/excelize/v2"
)

type adminEquity struct {
}

var AdminEquity = new(adminEquity)

func (s *adminEquity) Create(in model.CreateEquityActivityReq) (err error) {
	// 如果是白名单 校验白名单导入数据
	equityUserItems := []model.ImportItem{}
	if in.LimitType == model.EQUITY_ACTIVITY_LIMIT_TYPE2 {
		items, err2 := AdminEquity.HandelExcelUser(in)
		if err2 != nil {
			err = fmt.Errorf(err2.Error())
			return
		}
		if items.HaveErr {
			err = fmt.Errorf("表格数据存在异常数据，检查后重试")
			return
		}
		if items.SuccItems == nil {
			err = fmt.Errorf("导入数据为空，请重新输入")
			return
		}
		equityUserItems = items.SuccItems
	}
	// 校验导入数据结束

	// 获取详情
	equityItem, err2 := AdminEquity.Item(in.TemplateId)

	// 插入数据
	var tx *gdb.TX
	tx, err = g.DB().Begin()
	if err != nil {
		return
	}
	if err2 != nil {
		insterItem, insertEerr := tx.Model("equity_activity").Insert(&in)
		if insertEerr != nil {
			err = fmt.Errorf(insertEerr.Error())
			tx.Rollback()
			return
		}
		i, _ := insterItem.LastInsertId()
		equityItem.Id = int(i)
		_, err = tx.Model("activity").Insert(g.Map{
			"name":          in.Name,
			"start_time":    in.ActivityStartTime,
			"end_time":      in.ActivityEndTime,
			"publisher_id":  in.PublisherId,
			"activity_id":   equityItem.Id,
			"activity_type": model.ACTIVITY_TYPE_4,
		})
		if err != nil {
			tx.Rollback()
			return
		}
	} else {
		_, err = tx.Model("equity_activity").Where("template_id", in.TemplateId).Update(g.Map{
			"name":                in.Name,
			"price":               in.Price,
			"time_type":           in.TimeType,
			"activity_start_time": in.ActivityStartTime,
			"activity_end_time":   in.ActivityEndTime,
			"limit_type":          in.LimitType,
			"sub_limit_type":      in.SubLimitType,
			"limit_buy":           in.LimitBuy,
			"number":              in.Number,
			"status":              in.Status,
		})
		if err != nil {
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	if in.LimitType == model.EQUITY_ACTIVITY_LIMIT_TYPE2 {
		go AdminEquity.CreateEquityUser(in.PublisherId, equityItem.Id, equityUserItems)
	}
	return
}

// 解析处理用户数据
/*
	查询用户 用户存在 不用校验手机号是否符合规则
	判断是否有重复
	{
		"have_err" : false,是否有异常
		"items" : [ 数据集合
			{
				"phone":"",
				"user_id":""
				"err":""
			},
		],
		"total": 总条数
		"number": 总库存数
		"err_excel_url":"" 错误表格下载链接
	}
	如果有异常数据 写入 err_excel_url 返回表格链接
*/
func (s *adminEquity) HandelExcelUser(req model.CreateEquityActivityReq) (items model.ImportItems, err error) {

	//req.ExcelFile = "https://website-cdn.gfanx.com/developer/meta_world_id/2034061254251279543550bb04c840d30fae9f2ef282 (2).xlsx"

	if req.ExcelFile == "" {
		err = fmt.Errorf("请上传文件")
		return
	}
	_, bs, errs := gorequest.New().Get(req.ExcelFile).EndBytes()
	if len(errs) != 0 {
		err = fmt.Errorf("导入文件地址有误")
		return
	}
	f, e := excelize.OpenReader(bytes.NewReader(bs))
	if e != nil {
		err = fmt.Errorf("获取文件内容异常")
		return
	}
	rows, e := f.GetRows("Sheet1")
	if e != nil {
		err = fmt.Errorf("文件读取错误[解析手机号异常]")
		return
	}

	phoneItems, countItems := []string{}, []string{}

	for dx, v := range rows {
		if dx == 0 {
			continue
		}
		countItems = append(phoneItems, v[1])
		p := strings.TrimSpace(v[0])
		if p == "" {
			continue
		}
		phoneItems = append(phoneItems, p)
	}
	if len(phoneItems) <= 0 {
		err = fmt.Errorf("请输入手机号")
		return
	}
	if len(countItems) <= 0 {
		err = fmt.Errorf("请输入库存")
		return
	}
	if len(phoneItems) != len(countItems) {
		err = fmt.Errorf("手机号列数与库存列数不匹配")
		return
	}
	// 查询用户详情
	mobileString := strings.Join(phoneItems, ",")
	mobileString = strings.Replace(mobileString, " ", "", -1)
	phoneArr := strings.Split(mobileString, ",")
	userMap, e := provider.User.GetUserInfoByPhone(&map[string]interface{}{
		"phoneArr": phoneArr,
	})
	if e != nil {
		return items, e
	}

	haveErr := false
	total, number := 0, 0 //总条数 总库存数
	succItems, errItems := []model.ImportItem{}, []model.ImportItem{}
	for key, value := range rows {
		if key == 0 {
			continue
		}
		mobileRow := string(value[0])
		num, e := strconv.Atoi(value[1])
		number += num
		total += 1
		errMessage := ""
		if e != nil {
			errMessage = errMessage + "[" + e.Error() + "]"
		}
		if num <= 0 {
			errMessage = errMessage + "[数量异常]"
		}
		// 用户存在
		userItem := userMap[mobileRow]
		if userItem.UserId == "" {
			errMessage = errMessage + "[用户不存在]"
		}

		if errMessage != "" {
			haveErr = true
			errItem := model.ImportItem{
				ErrMessage: errMessage,
				LimitNum:   num,
				Phone:      mobileRow,
			}
			errItems = append(errItems, errItem)
		} else {
			succItem := model.ImportItem{
				UserId:   userItem.UserId,
				LimitNum: num,
				Phone:    mobileRow,
			}
			succItems = append(succItems, succItem)
		}
	}
	items.HaveErr = haveErr
	items.Total = total
	items.Number = number
	items.ErrItems = errItems
	items.SuccItems = succItems
	return
}

// 创建白名单用户数据
func (s *adminEquity) CreateEquityUser(PublishedId string, activityId int, equityUser []model.ImportItem) {
	for _, value := range equityUser {
		g.DB().Model("equity_user").Insert(g.Map{
			"publisher_id": PublishedId,
			"activity_id":  activityId,
			"user_id":      value.UserId,
			"phone":        value.Phone,
			"limit_num":    value.LimitNum,
		})
	}
}

// 获取详情
func (s *adminEquity) Item(templateId string) (ret model.EquityActivity, err error) {
	m := g.DB().Model("equity_activity")
	err = m.Where("template_id", templateId).Scan(&ret)
	if err != nil {
		return
	}
	return
}
