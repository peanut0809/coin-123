package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"meta_launchpad/model"
	"meta_launchpad/provider"
	"regexp"
	"strconv"
	"strings"
	"time"

	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	userModel "brq5j1d.gfanx.pro/meta_cloud/meta_service/app/user/model"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/parnurzeal/gorequest"
	"github.com/xuri/excelize/v2"
)

type adminEquity struct {
}

var AdminEquity = new(adminEquity)

func (s *adminEquity) Create(in model.CreateEquityActivityReq) (err error) {
	// 获取详情
	var equityItem *model.EquityActivity
	m := g.DB().Model("equity_activity")
	err = m.Where("template_id", in.TemplateId).Scan(&equityItem)
	if err != nil {
		err = fmt.Errorf("权益活动信息获取异常")
		return
	}
	if equityItem.Status == model.EQUITY_ACTIVITY_STATUS1 {
		err = fmt.Errorf("权益活动信息上架中，请勿重复上架")
		return
	}
	// 插入数据
	var tx *gdb.TX
	tx, err = g.DB().Begin()
	if err != nil {
		tx.Rollback()
		return
	}

	if equityItem == nil {
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
	}
	tx.Commit()
	// 如果是白名单 创建用户
	if in.LimitType == model.EQUITY_ACTIVITY_LIMIT_TYPE2 {
		go AdminEquity.CreateEquityUser(in.PublisherId, equityItem.Id, in)
	}
	return
}

/*
	查询用户 用户存在 不用校验手机号是否符合规则
	判断是否有重复
	{
		"have_err" : false,是否有异常
		"items" : [ 错误数据集合
			{
				"phone":"",
				"user_id":""
				"err":""
			},
		],
		"total": 总条数
		"number": 总库存数
		"errTotal":总错误条数
		"errNumber":总异常库存
	}
*/
func (s *adminEquity) HandelExcelUser(req model.CreateEquityActivityReq) (items model.ImportItems, err error) {

	// 获取总资产数
	// UsedAssets := 0
	// ai, err2 := AdminEquity.GetCanUsedAssetsCountByTemplate(req.AppId, req.TemplateId)
	// if err2 != nil {
	// 	err = fmt.Errorf("获取资产库存异常")
	// 	return
	// }
	// UsedAssets = ai.Count

	if req.ExcelFile == "" {
		err = fmt.Errorf("请上传文件")
		g.Log().Error(req.PublisherId + "导入白名单异常!" + "请上传文件")
		return
	}
	_, bs, errs := gorequest.New().Get(req.ExcelFile).EndBytes()
	if len(errs) != 0 {
		err = fmt.Errorf("导入文件地址有误")
		g.Log().Error(req.PublisherId + "导入白名单异常!" + "导入文件地址有误" + "[" + req.ExcelFile + "]")
		return
	}
	f, e := excelize.OpenReader(bytes.NewReader(bs))
	if e != nil {
		err = fmt.Errorf("获取文件内容异常")
		g.Log().Error(req.PublisherId + "导入白名单异常!获取文件内容异常" + e.Error())
		return
	}
	rows, e := f.GetRows("Sheet1")
	if e != nil {
		err = fmt.Errorf("文件读取错误[解析Sheet1手机/数量异常]")
		g.Log().Error(req.PublisherId + "导入白名单异常!文件读取错误[解析Sheet1手机/数量异常]" + e.Error())
		return
	}

	phoneItems, _ := []string{}, []string{}
	for dx, v := range rows {
		if dx == 0 {
			continue
		}
		if len(v) < 2 {
			err = fmt.Errorf("列文件数据错误[两列长度不匹配]")
			return
		}
		n := strings.TrimSpace(v[1])
		if n == "" {
			continue
		}
		//countItems = append(phoneItems, v[1])
		p := strings.TrimSpace(v[0])
		if p == "" {
			continue
		}
		phoneItems = append(phoneItems, p)
	}

	// 查询用户详情
	var userMap map[string]provider.GetUserInfoRet
	var equityUserMap map[string]model.EquityUser

	// 如果是创建 校验用户信息
	if req.IsCreate {
		userMap, equityUserMap, err = AdminEquity.HandelUserItems(phoneItems)
		if err != nil {
			return items, err
		}
	}

	/*
		rows 表格数据
		userMap rpc获取用户数据集合
		equityUserMap 白名单表用户是否存在
		req.IsCreate 如果是创建 校验用户+创建用户
	*/
	haveErr, number, succItems, errItems := AdminEquity.HandelExcelRowErr(rows, userMap, equityUserMap, req.IsCreate)

	items.HaveErr = haveErr
	items.Total = len(rows) - 1
	items.Number = number
	items.ErrItems = errItems
	items.SuccItems = succItems
	//items.AssetsCount = UsedAssets
	return
}

// 创建白名单用户数据
func (s *adminEquity) CreateEquityUser(PublishedId string, activityId int, in model.CreateEquityActivityReq) {
	equityUserItems := []model.ImportItem{}
	if in.LimitType == model.EQUITY_ACTIVITY_LIMIT_TYPE2 {
		in.IsCreate = true
		items, err2 := AdminEquity.HandelExcelUser(in)
		if err2 != nil {
			g.Log().Error(PublishedId + "导入白名单异常!" + err2.Error())
			return
		}
		if items.HaveErr {
			g.Log().Error(PublishedId + "导入白名单异常!表格数据存在异常数据，检查后重试")
			return
		}
		if items.SuccItems == nil {
			g.Log().Error(PublishedId + "导入白名单异常!导入数据为空，请重新输入")
			return
		}
		equityUserItems = append(items.SuccItems, items.ErrItems...)
		in.Number = items.Number
	}
	// 校验导入数据结束
	number := 0

	var tx *gdb.TX
	tx, err := g.DB().Begin()
	if err != nil {
		tx.Rollback()
		return
	}
	for _, value := range equityUserItems {
		_, err := tx.Model("equity_user").Insert(g.Map{
			"publisher_id": PublishedId,
			"activity_id":  activityId,
			"user_id":      value.UserId,
			"phone":        value.Phone,
			"limit_num":    value.LimitNum,
		})
		if err != nil {
			tx.Rollback()
			return
		}
		number += value.LimitNum
	}
	_, err = tx.Model("equity_activity").Where(g.Map{
		"publisher_id": PublishedId,
		"id":           activityId,
	}).Update(g.Map{
		"number":       number,
		"total_number": number,
	})
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
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

// 下架活动
func (s *adminEquity) Invalid(EquityId int) (err error) {
	m := g.DB().Model("equity_activity")
	_, err = m.Where("id", EquityId).Update(g.Map{
		"status": model.EQUITY_ACTIVITY_STATUS2,
	})
	if err != nil {
		return
	}
	return
}

func (s *adminEquity) AssetsCount(req model.CreateEquityActivityReq) (count map[string]int, err error) {
	ai, err2 := AdminEquity.GetCanUsedAssetsCountByTemplate(req.AppId, req.TemplateId)
	if err2 != nil {
		err = fmt.Errorf(err2.Error())
		return
	}
	count = make(map[string]int)
	count["AssetsCount"] = ai.Count
	return
}

// 用户明细
func (s *adminEquity) UserItems(in model.EquityUserReq) (list model.EquityUserFull, err error) {
	m := g.DB().Model("equity_user").Where("publisher_id", in.PublisherId).Where("activity_id", in.EquityId)
	if in.Phone > 0 {
		m = m.Where("phone", in.Phone)
	}
	if in.Status > 0 {
		m = m.Where("status", in.Status)
	}
	total, err := m.Count()
	if err != nil {
		err = gerror.New("获取总行数失败")
		return
	}
	list.Total = total
	userList := make([]model.EquityUser, 0)
	err = m.Order("id DESC").Page(in.Page, in.PageSize).Scan(&userList)
	if err != nil {
		return
	}
	list.List = userList
	if err != nil {
		return
	}
	return
}

// 权益活动记录
func (s *adminEquity) EquityActivityItems(in model.AdminEquityReq) (ret model.EquityActivityList, err error) {
	m := g.DB().Model("equity_activity").Where("publisher_id", in.PublisherId)
	if in.Status > 0 {
		m = m.Where("status", in.Status)
	}
	if in.TemplateId != "" {
		m = m.Where("template_id", in.TemplateId)
	}
	if in.StartDate != "" {
		m = m.Where("created_at >= ", in.StartDate)
	}

	if in.EndDate != "" {
		m = m.Where("created_at <= ", in.EndDate)
	}

	if in.Name != "" {
		m = m.WhereLike("name", "%"+in.Name+"%")
	}

	total, err := m.Count()
	if err != nil {
		err = gerror.New("获取总行数失败")
		return
	}
	ret.Total = total
	rs := make([]*model.EquityActivity, 0)
	err = m.Order("id DESC").Page(in.Page, in.PageSize).Scan(&rs)
	if err != nil {
		return
	}
	ret.List = rs
	if err != nil {
		return
	}
	return
}

func (s *adminEquity) OrderItems(in model.AdminEquityOrderReq) (ret model.EquityOrderList, err error) {
	m := g.DB().Model("equity_orders").Where("publisher_id", in.PublisherId)
	if in.Phone > 0 {
		m = m.Where("phone", in.Phone)
	}
	if in.Status > 0 {
		m = m.Where("status", in.Status)
	}
	if in.StartDate != "" {
		m = m.Where("created_at < ", in.StartDate)
	}
	if in.EndDate != "" {
		m = m.Where("created_at > ", in.EndDate)
	}
	if in.MinPrice > 0 {
		m = m.Where("real_fee > ", in.MinPrice)
	}
	if in.MaxPrice > 0 {
		m = m.Where("real_fee < ", in.MaxPrice)
	}
	if in.OrderNo != "" {
		m = m.Where("order_no ", in.OrderNo)
	}
	total, err := m.Count()
	if err != nil {
		err = gerror.New("获取总行数失败")
		return
	}
	ret.Total = total

	var items []*model.EquityOrder
	err = m.Order("id DESC").Page(in.Page, in.PageSize).Scan(&items)
	if err != nil {
		return
	}
	for _, v := range items {
		lastSec := v.PayExpireAt.Unix() - time.Now().Unix()
		if lastSec <= 0 {
			lastSec = 0
		}
		ret.List = append(ret.List, &model.EquityOrderFull{
			EquityOrder: v,
			PriceYuan:   fmt.Sprintf("%.2f", float64(v.Price)/100),
			RealFeeYuan: fmt.Sprintf("%.2f", float64(v.RealFee)/100),
			LastSec:     lastSec,
		})
	}
	return
}

// 定时任务到期活动下架
func (s *adminEquity) OffShelvesEquityActivity(ctx context.Context, req interface{}, result *interface{}) (err error) {
	var items []*model.EquityActivity
	timeCrv := time.Now().Unix()
	nowTime := time.Unix(timeCrv, 0).Format("2006-01-02 15:04:05")
	m := g.DB().Model("equity_activity").Where("activity_end_time", nowTime)
	m.Scan(&items)
	if items != nil {
		g.DB().Model("equity_activity").Where("activity_end_time", nowTime).Limit(len(items)).Update(g.Map{
			"status": model.EQUITY_ACTIVITY_STATUS2,
		})
		return
	}
	return
}

// 获取导入表格用户信息
func (s *adminEquity) HandelUserItems(phoneItems []string) (userMap map[string]provider.GetUserInfoRet, equityUserMap map[string]model.EquityUser, err error) {
	mobileString := strings.Join(phoneItems, ",")
	mobileString = strings.Replace(mobileString, " ", "", -1)
	phoneArr := strings.Split(mobileString, ",")
	userMap, err = provider.User.GetUserInfoByPhone(&map[string]interface{}{
		"phoneArr": phoneArr,
	})
	if err != nil {
		g.Log().Error("导入白名单异常!rpc获取用户数据异常" + err.Error())
		return
	}

	// 获取activity_id关联手机号用户信息
	var users []model.EquityUser
	err = g.DB().Model("equity_user").Where("status", 1).Where("phone IN (?)", phoneArr).Scan(&users)
	if err != nil {
		g.Log().Error("导入白名单异常!权益活动id获取用户异常" + err.Error())
		return
	}
	equityTempUserMap := make(map[string]model.EquityUser)
	for _, v := range users {
		equityTempUserMap[v.Phone] = v
	}
	return userMap, equityTempUserMap, err
}

func (s *adminEquity) HandelExcelRowErr(rows [][]string, userMap map[string]provider.GetUserInfoRet, equityUserMap map[string]model.EquityUser, isCreate bool) (haveErr bool, number int, succItems, errItems []model.ImportItem) {
	m := make(map[interface{}]interface{})
	for key, value := range rows {
		if key == 0 {
			continue
		}
		mobileRow := string(value[0])
		num, _ := strconv.Atoi(value[1])

		// 导入用户总购买数量
		number += num

		errMessage := ""
		result, _ := regexp.MatchString(`^(1[3|4|5|6|7|8|9][0-9]\d{4,8})$`, mobileRow)
		if !result {
			errMessage = errMessage + "[手机号格式异常]"
		}

		_, ok := m[mobileRow]
		if ok {
			errMessage = errMessage + "[手机号重复]"
		} else {
			m[mobileRow] = value
		}
		if num <= 0 {
			errMessage = errMessage + "[数量异常]"
		}
		userItem := userMap[mobileRow]

		//如果是创建校验用户+创建用户
		if isCreate {
			if userItem.UserId == "" {
				creatUserInfo, createUserErr := AdminEquity.UsersRegistUserByPhone(mobileRow)
				if createUserErr != nil {
					errMessage = errMessage + createUserErr.Error()
				} else {
					userItem.UserId = creatUserInfo.UserId
				}
			}
			// 用户在白名单已经存在
			equityUserItem := equityUserMap[mobileRow]
			if equityUserItem.Phone != "" {
				errMessage = errMessage + "[用户已经在白名单]"
			}
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
	return
}

// rpc查询资产库存
func (s *adminEquity) GetCanUsedAssetsCountByTemplate(appId, templateId string) (ret *model.AssetItems, err error) {

	userAsset := &model.AssetItems{}
	params := &map[string]interface{}{
		"appId":      appId,
		"templateId": templateId,
	}
	result, err := utils.SendJsonRpc(context.Background(), "asset", "Asset.GetCanUsedAssetsCountByTemplate", params)

	if err != nil {
		g.Log().Errorf("GetCanUsedAssetsCountByTemplate err:%v", err)
		return
	}
	err = json.Unmarshal([]byte(gconv.String(result)), userAsset)
	if err != nil {
		g.Log().Error(err)
		return nil, err
	}
	return userAsset, nil
}

// rpc创建注册用户
func (s *adminEquity) UsersRegistUserByPhone(phone string) (ret *userModel.Users, err error) {
	// 用户信息
	user := &userModel.Users{}
	params := &map[string]interface{}{
		"phone": phone,
	}
	// 获取用户信息
	result, err := utils.SendJsonRpc(context.Background(), "ucenter", "Users.UsersRegistUserByPhone", params)
	if err != nil {
		g.Log().Error(err)
		return
	}
	err = json.Unmarshal([]byte(gconv.String(result)), user)
	if err != nil {
		g.Log().Error(err)
		return
	}
	return user, nil
}
