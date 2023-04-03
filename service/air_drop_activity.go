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

	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	"github.com/gogf/gf/frame/g"
	"github.com/parnurzeal/gorequest"
	"github.com/xuri/excelize/v2"
)

type airDropActivity struct {
}

var AirDropActivity = new(airDropActivity)

func (s *airDropActivity) Items(req *model.AirDropActivityItemReq) (ret model.AirDropActivityItemsRex, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	sql := g.DB().Model("air_drop_activity")
	if req.Name != "" {
		sql = sql.Where("name like ? ", "%"+req.Name+"%")
	}
	if req.Type != "" {
		sql = sql.Where("type", req.Type)
	}
	i, err := sql.Count()
	ret.Total = i
	if err != nil {
		return
	}
	err = sql.OrderDesc("id").Page(req.Page, req.PageSize).Scan(&ret.List)
	if err != nil {
		return
	}
	for _, ada := range ret.List {
		successCount, err := g.DB().Model("air_drop_activity_record").Where("activity_id", ada.Id).Where("status", model.Air_Drop_Status_01).Count()
		if err == nil {
			ada.SuccessCount = successCount
		}
		errCount, err := g.DB().Model("air_drop_activity_record").Where("activity_id", ada.Id).Where("status", model.Air_Drop_Status_02).Count()
		if err == nil {
			ada.ErrorCount = errCount
		}
	}
	return
}

func (s *airDropActivity) Item(req *model.AirDropActivityItemReq) (ret model.AirDropActivityItemRex, err error) {
	if req.DropId <= 0 {
		err = fmt.Errorf("drop id err")
		return
	}

	var airDropActivity *model.AirDropActivity
	err = g.DB().Model(&model.AirDropActivity{}).Where("id", req.DropId).Scan(&airDropActivity)
	if err != nil {
		return
	}
	if airDropActivity == nil {
		err = fmt.Errorf("airDropActivity not exit")
		return
	}

	successCount, err := g.DB().Model("air_drop_activity_record").Where("activity_id", req.DropId).Where("status", model.Air_Drop_Status_01).Count()
	if err == nil {
		airDropActivity.SuccessCount = successCount
	}
	errCount, err := g.DB().Model("air_drop_activity_record").Where("activity_id", req.DropId).Where("status", model.Air_Drop_Status_02).Count()
	if err == nil {
		airDropActivity.ErrorCount = errCount
	}
	ret.AirDropActivity = airDropActivity

	recordSql := g.DB().Model(&model.AirDropActivityRecord{}).Where("activity_id", req.DropId)
	if req.Phone != "" {
		recordSql = recordSql.Where("phone", req.Phone)
	}
	if req.Status > 0 {
		recordSql = recordSql.Where("status", req.Status)
	}
	ret.Total, err = recordSql.Count()
	if err != nil {
		return
	}
	err = recordSql.OrderDesc("id").Page(req.Page, req.PageSize).Scan(&ret.Items)
	if err != nil {
		return
	}

	return
}
func (s *airDropActivity) AirDrop(req *model.AirDropActivityReq) (ret interface{}, err error) {
	if req.ExcelFile == "" {
		err = fmt.Errorf("excelFile err")
		return
	}
	if req.Name == "" || req.Remark == "" {
		err = fmt.Errorf("name or remark err")
		return
	}
	typeItem := model.Air_Drop_Type_Map
	_, ok := typeItem[req.Type]
	if !ok {
		err = fmt.Errorf("type should in [speed,crystal]")
		return
	}
	mobileItems, err := s.MakeExcelItems(req.ExcelFile)
	if err != nil {
		return
	}
	// 创建主活动
	config, _ := json.Marshal(req)
	item := model.AirDropActivity{
		OrderNo: utils.Generate(),
		Config:  string(config),
		Type:    req.Type,
		Name:    req.Name,
		Remark:  req.Remark,
	}
	r, err := g.DB().Model(&model.AirDropActivity{}).Insert(&item)
	if err != nil {
		return
	}
	dropId, err := r.LastInsertId()
	if err != nil {
		return
	}
	item.Id = int(dropId)
	if req.Type == model.Air_Drop_Type_Speed {
		go s.MakeSpeedDrop(item, mobileItems)
	}
	if req.Type == model.Air_Drop_Type_Crystal {
		go s.MakeCrystalDrop(item, mobileItems)
	}
	return
}

// 加速次数空投
func (s *airDropActivity) MakeSpeedDrop(req model.AirDropActivity, mobileItems []model.MobileCollect) {
	// 手机号集合 手机号 用户id 数量 message(创建用户、手机号格式异常)
	mc := s.MakeMobileItems(mobileItems)
	for _, mobileInfo := range mc {
		errMessage := ""
		mobileInfo.ActivityId = req.Id
		err := s.MarktingUserSpeedNum(mobileInfo)
		if err != nil {
			errMessage = errMessage + err.Error()
		}
		if errMessage != "" {
			mobileInfo.Message = mobileInfo.Message + "【" + errMessage + "】"
		}
		status := model.Air_Drop_Status_01
		if mobileInfo.HaveErr == true || errMessage != "" {
			status = model.Air_Drop_Status_02
		}
		_, err = g.DB().Model("air_drop_activity_record").Insert(&model.AirDropActivityRecord{
			UserId:       mobileInfo.UserId,
			Phone:        mobileInfo.Mobile,
			ActivityId:   req.Id,
			Number:       mobileInfo.Number,
			ActivityType: req.Type,
			Status:       status,
			Message:      errMessage,
		})
		if err != nil {
			g.Log().Info("MakeSpeedDrop err " + mobileInfo.UserId + err.Error())
			continue
		}
	}
}

// 元晶空投
func (s *airDropActivity) MakeCrystalDrop(req model.AirDropActivity, mobileItems []model.MobileCollect) {
	mc := s.MakeMobileItems(mobileItems)
	for _, mobileInfo := range mc {
		errMessage := ""
		//发放元晶
		xJTransferReq := &map[string]interface{}{
			"userId":   mobileInfo.UserId,
			"category": 1,
			"amount":   mobileInfo.Number,
			"source":   36,
			"orderNo":  req.OrderNo,
			//"orderNo":  utils.Generate(),
		}
		err := provider.User.YJTransfer(xJTransferReq)
		if err != nil {
			errMessage = errMessage + err.Error()
		}
		if errMessage != "" {
			mobileInfo.Message = mobileInfo.Message + "【" + errMessage + "】"
		}
		status := model.Air_Drop_Status_01
		if mobileInfo.HaveErr == true || errMessage != "" {
			status = model.Air_Drop_Status_02
		}
		_, err = g.DB().Model("air_drop_activity_record").Insert(&model.AirDropActivityRecord{
			UserId:       mobileInfo.UserId,
			Phone:        mobileInfo.Mobile,
			ActivityId:   req.Id,
			Number:       mobileInfo.Number,
			ActivityType: req.Type,
			Status:       status,
			Message:      errMessage,
		})
		if err != nil {
			g.Log().Info("MakeCrystalDrop err " + mobileInfo.UserId + err.Error())
			continue
		}
	}

}

// 获取用户 & 创建用户
func (s *airDropActivity) MakeMobileItems(mobileItems []model.MobileCollect) (mobileCollect []model.MobileCollect) {
	phoneItems := make([]string, 0)
	for _, mc := range mobileItems {
		phoneItems = append(phoneItems, mc.Mobile)
	}
	if len(phoneItems) <= 0 {
		g.Log().Info("MakeMobileItems phoneItems lens <= 0")
		return
	}
	mobileString := strings.Join(phoneItems, ",")
	phoneArr := strings.Split(mobileString, ",")
	userMap, err := provider.User.GetUserInfoByPhone(&map[string]interface{}{
		"phoneArr": phoneArr,
	})
	if err != nil {
		g.Log().Error("MakeMobileItems rpc get user err" + err.Error())
		return
	}
	for _, mobilItem := range mobileItems {
		result, err := regexp.MatchString("^1[3|4|5|6|7|8|9]{1}\\d{9}$", mobilItem.Mobile)
		if !result {
			mobilItem.HaveErr = true
			mobilItem.Message = "手机号格式异常" + err.Error()
			continue
		}
		userInfo := userMap[mobilItem.Mobile]
		if userInfo.Phone == "" {
			guir, err := AdminEquity.UsersRegistUserByPhone(mobilItem.Mobile)
			if err != nil {
				mobilItem.HaveErr = true
				mobilItem.Message = "UsersRegistUserByPhone err" + err.Error()
				continue
			} else {
				mobilItem.Message = "(注册新用户)"
				mobilItem.UserId = guir.UserId
			}
		} else {
			mobilItem.UserId = userInfo.UserId
		}
		mobileCollectItem := model.MobileCollect{
			Mobile:  mobilItem.Mobile,
			Number:  mobilItem.Number,
			UserId:  mobilItem.UserId,
			Message: mobilItem.Message,
		}
		mobileCollect = append(mobileCollect, mobileCollectItem)
	}
	return
}

// 解析表
func (s *airDropActivity) MakeExcelItems(excelFile string) (mobileItems []model.MobileCollect, err error) {
	_, bs, errs := gorequest.New().Get(excelFile).EndBytes()
	if len(errs) != 0 {
		err = fmt.Errorf("excel file address err")
		return
	}
	f, e := excelize.OpenReader(bytes.NewReader(bs))
	if e != nil {
		err = fmt.Errorf("excel file read err")
		return
	}
	rows, e := f.GetRows("Sheet1")
	if e != nil {
		err = fmt.Errorf("excel file sheet1 err")
		return
	}
	for dx, v := range rows {
		if len(v) < 2 {
			err = fmt.Errorf("列文件数据错误[两列长度不匹配]")
			return
		}
		if dx == 0 {
			continue
		}
		mobile := strings.TrimSpace(v[0])
		number := strings.TrimSpace(v[1])
		speedNum, _ := strconv.Atoi(number)
		info := model.MobileCollect{
			Mobile: mobile,
			Number: speedNum,
		}
		mobileItems = append(mobileItems, info)
	}
	if len(mobileItems) <= 0 {
		err = fmt.Errorf("mobile items len err")
		return
	}
	return
}

func (s *airDropActivity) MarktingUserSpeedNum(item model.MobileCollect) (err error) {
	// 获取用户信息
	params := g.Map{
		"userId":     item.UserId,
		"number":     item.Number,
		"from":       2,
		"activityId": item.ActivityId,
	}
	_, err = utils.SendJsonRpc(context.Background(), "activity", "FormulaMarkting.MarktingUserSpeedNum", params)
	if err != nil {
		g.Log().Error(err)
		return
	}
	return
}
