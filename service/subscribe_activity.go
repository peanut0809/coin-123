package service

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"meta_launchpad/cache"
	"meta_launchpad/model"
	"meta_launchpad/provider"
	"time"
)

const SubSetResultKey = "meta_launchpad:activity_result:%s"

type subscribeActivity struct {
}

var SubscribeActivity = new(subscribeActivity)

//func (s *subscribeActivity) GetList() (ret map[string][]model.SubscribeActivityFull, err error) {
//	ret = make(map[string][]model.SubscribeActivityFull)
//	var as []model.SubscribeActivity
//	m := g.DB().Model("subscribe_activity")
//	now := time.Now()
//	err = m.Where("publisher_id != 'ZHW' AND start_time <= ? AND pay_end_time >= ?", now, now).Order("activity_start_time ASC").Scan(&as)
//	if err != nil {
//		return
//	}
//	ret["priority"] = make([]model.SubscribeActivityFull, 0)
//	ret["general"] = make([]model.SubscribeActivityFull, 0)
//	for _, v := range as {
//		var item model.SubscribeActivityFull
//		item.Name = v.Name
//		item.CoverImgUrl = v.CoverImgUrl
//		item.SumNum = v.SumNum
//		item.Alias = v.Alias
//		item.PriceYuan = fmt.Sprintf("%.2f", float64(v.Price)/100)
//		item.Status = s.GetActivityStatusV2(v)
//		if v.ActivityType == 1 {
//			ret["priority"] = append(ret["priority"], item)
//		} else {
//			ret["general"] = append(ret["general"], item)
//		}
//	}
//	return
//}

//func (s *subscribeActivity) GetActivityStatus(v model.SubscribeActivity) (status int, lastSec int64) {
//	now := time.Now()
//	if now.Unix() <= v.ActivityStartTime.Unix() { //距开始
//		status = model.STATUS_AWAY_START
//		lastSec = v.ActivityStartTime.Unix() - now.Unix()
//	} else if now.Unix() >= v.ActivityStartTime.Unix() && now.Unix() < v.ActivityEndTime.Unix() { //距结束
//		status = model.STATUS_AWAY_END
//		lastSec = v.ActivityEndTime.Unix() - now.Unix()
//	} else if now.Unix() >= v.ActivityEndTime.Unix() && now.Unix() <= v.OpenAwardTime.Unix() { //距开奖
//		status = model.STATUS_AWAY_AWARD
//		lastSec = v.OpenAwardTime.Unix() - now.Unix()
//	} else if now.Unix() > v.OpenAwardTime.Unix() && now.Unix() <= v.PayEndTime.Unix() { //距付款截止
//		status = model.STATUS_AWAY_PAY_END
//		lastSec = v.PayEndTime.Unix() - now.Unix()
//	}
//	return
//}

func (s *subscribeActivity) GetActivityStatusV2(v model.SubscribeActivity, userId string) (status int) {
	now := time.Now()
	if now.Unix() <= v.ActivityStartTime.Unix() { //未开始
		status = model.STATUS_AWAY_START
	} else if now.Unix() >= v.ActivityStartTime.Unix() && now.Unix() < v.ActivityEndTime.Unix() { //距结束
		status = model.STATUS_ING
		record, _ := SubscribeRecord.GetSubscribeRecords(v.Id, userId)
		if record != nil {
			status = model.STATUS_SUBED
		}
		//ret.Subed = record != nil
	} else if now.Unix() >= v.ActivityEndTime.Unix() && now.Unix() <= v.OpenAwardTime.Unix() { //待公布
		status = model.STATUS_AWAIT_OPEN
	} else if now.Unix() > v.OpenAwardTime.Unix() && now.Unix() <= v.PayEndTime.Unix() { //待付款
		//status = model.STATUS_AWAIT_PAY
		status = model.STATUS_END
	} else if now.Unix() > v.PayEndTime.Unix() { //已结束
		status = model.STATUS_END
	}
	return
}

func (s *subscribeActivity) GetValidDetail(alias string) (ret *model.SubscribeActivity, err error) {
	m := g.DB().Model("subscribe_activity")
	now := time.Now()
	err = m.Where("start_time <= ? AND alias = ?", now, alias).Scan(&ret)
	if err != nil {
		return
	}
	if ret == nil {
		err = fmt.Errorf("活动不存在")
		return
	}
	return
}

func (s *subscribeActivity) GetSimpleDetail(id int) (ret model.SubscribeActivity, err error) {
	err = g.DB().Model("subscribe_activity").Where("id = ?", id).Scan(&ret)
	return
}

func (s *subscribeActivity) GetSimpleDetailByAlias(alias string) (ret *model.SubscribeActivity, err error) {
	err = g.DB().Model("subscribe_activity").Where("alias = ?", alias).Scan(&ret)
	return
}

func (s *subscribeActivity) GetListSimple(ids []int) (ret []model.SubscribeActivity, err error) {
	err = g.DB().Model("subscribe_activity").Where("id in (?)", ids).Scan(&ret)
	return
}

func (s *subscribeActivity) GetDetail(alias, userId, publisherId string) (ret model.SubscribeActivityFull, err error) {
	as, e := s.GetValidDetail(alias)
	if e != nil {
		err = e
		return
	}
	ret.Name = as.Name
	ret.CoverImgUrl = as.CoverImgUrl
	ret.SumNum = as.SumNum
	ret.Alias = as.Alias
	ret.PriceYuan = fmt.Sprintf("%.2f", float64(as.Price)/100)
	ret.Status = s.GetActivityStatusV2(*as, userId)
	ret.ActivityType = as.ActivityType
	ret.AssetIntro = as.AssetIntro
	ret.ActivityIntro = as.ActivityIntro
	ret.SubSumPeople = as.SubSumPeople
	ret.SubSum = as.SubSum
	ret.Steps = append(ret.Steps, model.SubscribeActivityFullStep{
		Txt:     "开放抽签",
		TimeStr: as.ActivityStartTime.Layout("01-02 15:04"),
	})
	ret.Steps = append(ret.Steps, model.SubscribeActivityFullStep{
		Txt:     "抽签结束",
		TimeStr: as.ActivityEndTime.Layout("01-02 15:04"),
	})
	ret.Steps = append(ret.Steps, model.SubscribeActivityFullStep{
		Txt:     "公布时间",
		TimeStr: as.OpenAwardTime.Layout("01-02 15:04"),
	})
	ret.Steps = append(ret.Steps, model.SubscribeActivityFullStep{
		Txt:     "付款截止",
		TimeStr: as.PayEndTime.Layout("01-02 15:04"),
	})

	//获取应用信息
	appInfo, _ := provider.Developer.GetAppInfo(as.AppId)
	ret.AssetCateString = appInfo.Data.CnName
	//获取资产分类
	templateInfo, _ := provider.Developer.GetAssetsTemplate(as.AppId, as.TemplateId)
	for _, v := range templateInfo.CateList {
		ret.AssetCateString += fmt.Sprintf("-%s", v.CnName)
	}
	assetDetail, e := provider.Asset.GetMateDataByAm(&map[string]interface{}{
		"appId":      as.AppId,
		"templateId": as.TemplateId,
	})
	if e != nil {
		g.Log().Errorf("资产详情错误：%v", e)
		err = fmt.Errorf("获取资产信息失败")
		return
	}
	publisherInfo, e := provider.Developer.GetPublishInfo(publisherId)
	if e != nil {
		g.Log().Errorf("发行商异常：%v，%s", e, publisherId)
		err = fmt.Errorf("获取发行商失败")
		return
	}
	ret.AssetPic = assetDetail.AssetPic
	ret.ChainName = publisherInfo.ChainName
	ret.ChainAddr = publisherInfo.ChainAddr
	ret.ChainType = publisherInfo.ChainType
	ret.AssetTotal = assetDetail.Total
	ret.AssetCreateAt = assetDetail.CreateTime
	ret.AssetDetailImg = templateInfo.DetailImg
	ret.NfrDay = as.NfrSec / 3600 / 24
	return
}

func (s *subscribeActivity) GetAssetMaxBuyNum(aid int, userId string) (num int, err error) {
	conditions, e := SubscribeCondition.GetList(aid)
	if e != nil {
		err = e
		return
	}
	sumAsset, e := SubscribeCondition.GetConditionsAsset(userId, conditions)
	if e != nil {
		err = e
		return
	}
	for _, condition := range conditions {
		var assets []model.Asset
		assets, err = SubscribeCondition.GetOneConditionAsset(condition, sumAsset)
		if err != nil {
			return
		}
		num += len(assets) * condition.BuyNum
	}
	return
}

func (s *subscribeActivity) GetValidTicketInfo(ticketInfoStr string) (ticketInfo []model.TicketInfoJson, err error) {
	var sqlTicketInfo []model.TicketInfoJson
	err = json.Unmarshal([]byte(ticketInfoStr), &sqlTicketInfo)
	if err != nil {
		return
	}
	for _, v := range sqlTicketInfo {
		if v.Use {
			ticketInfo = append(ticketInfo, v)
		}
	}
	return
}

func (s *subscribeActivity) GetMaxBuyNum(alias string, userId string) (ticketInfo []model.TicketInfoJson, as *model.SubscribeActivity, err error) {
	as, err = s.GetValidDetail(alias)
	if err != nil {
		return
	}
	//获取用户元晶余额
	_, userInfoMap, _ := provider.User.GetUserInfo([]string{userId})
	if len(userInfoMap) == 0 {
		err = fmt.Errorf("用户信息异常")
		return
	}
	//获取用户月票
	//monthTicketInfo, e := provider.User.GetUserMonthTicket(userId)
	//if e != nil {
	//	err = e
	//	return
	//}
	//if monthTicketInfo == nil {
	//	err = fmt.Errorf("月票数据异常")
	//	return
	//}
	ticketInfo, err = s.GetValidTicketInfo(as.TicketInfo)
	if err != nil {
		return
	}
	for k, v := range ticketInfo {
		if v.Type == model.TICKET_CRYSTAL {
			ticketInfo[k].Num = userInfoMap[userId].Crystal
		}
		if v.Type == model.TICKET_MONTH {
			//ticketInfo[k].Num = monthTicketInfo.MonthTicket
		}
	}
	if as.ActivityType == 1 { //优先购
		assetNum, e := s.GetAssetMaxBuyNum(as.Id, userId)
		if e != nil {
			err = e
			return
		}
		for k, v := range ticketInfo {
			if v.Type == model.TICKET_MONEY {
				ticketInfo[k].MaxBuyNum = assetNum
			} else if v.Type == model.TICKET_CRYSTAL {
				if v.UnitNum != 0 {
					ticketInfo[k].MaxBuyNum = userInfoMap[userId].Crystal / v.UnitNum
					if ticketInfo[k].MaxBuyNum >= assetNum {
						ticketInfo[k].MaxBuyNum = assetNum
					}
				} else {
					ticketInfo[k].MaxBuyNum = assetNum
				}
			} else if v.Type == model.TICKET_MONTH {
				if v.UnitNum != 0 {
					//ticketInfo[k].MaxBuyNum = monthTicketInfo.MonthTicket / v.UnitNum
					if ticketInfo[k].MaxBuyNum >= assetNum {
						ticketInfo[k].MaxBuyNum = assetNum
					}
				}
			}
		}
	}
	if as.ActivityType == 2 { //普通购
		wallerRet, _ := provider.Wallet.WalletAuthenticationState(userId)

		for k, v := range ticketInfo {
			if v.Type == model.TICKET_MONEY {
				ticketInfo[k].MaxBuyNum = as.GeneralBuyNum
				if as.GeneralNumMethod == 1 {
					if wallerRet != nil {
						ticketInfo[k].MaxBuyNum = wallerRet.Account / as.Price
					}
				}
			} else if v.Type == model.TICKET_CRYSTAL {
				if v.UnitNum != 0 {
					ticketInfo[k].MaxBuyNum = userInfoMap[userId].Crystal / v.UnitNum
					if ticketInfo[k].MaxBuyNum >= as.GeneralBuyNum {
						ticketInfo[k].MaxBuyNum = as.GeneralBuyNum
					}
				} else {
					ticketInfo[k].MaxBuyNum = as.GeneralBuyNum
					if as.GeneralNumMethod == 1 {
						if wallerRet != nil {
							ticketInfo[k].MaxBuyNum = wallerRet.Account / as.Price
						}
					}
				}
			} else if v.Type == model.TICKET_MONTH {
				if v.UnitNum != 0 {
					//ticketInfo[k].MaxBuyNum = monthTicketInfo.MonthTicket / v.UnitNum
					if ticketInfo[k].MaxBuyNum >= as.GeneralBuyNum {
						ticketInfo[k].MaxBuyNum = as.GeneralBuyNum
					}
				} else {
					ticketInfo[k].MaxBuyNum = as.GeneralBuyNum
					if as.GeneralNumMethod == 1 {
						if wallerRet != nil {
							ticketInfo[k].MaxBuyNum = wallerRet.Account / as.Price
						}
					}
				}
			}
		}
	}
	return
}

func (s *subscribeActivity) DoSubVerify(in model.DoSubReq) (oneTicketInfo model.TicketInfoJson, activityInfo *model.SubscribeActivity, err error) {
	ticketInfo, as, e := s.GetMaxBuyNum(in.Alias, in.UserId)
	if e != nil {
		err = e
		return
	}
	if as.Disable == 1 {
		err = fmt.Errorf("活动已被禁用")
		return
	}
	now := time.Now()
	if !(now.Unix() > as.ActivityStartTime.Unix() && now.Unix() < as.ActivityEndTime.Unix()) {
		err = fmt.Errorf("认购时间已结束")
		return
	}
	//检查超时行为
	//gv, e := g.Redis().DoVar("GET", fmt.Sprintf(cache.SUB_PAY_TIMEOUT, in.UserId, as.ActivityType))
	//if e != nil {
	//	err = e
	//	return
	//}
	//if !gv.IsEmpty() {
	//	err = fmt.Errorf("已超时%s次未付款，资格已被取消", gv.String())
	//	return
	//}
	activityInfo = as
	for _, v := range ticketInfo {
		if v.Type == in.Type {
			oneTicketInfo = v
			break
		}
	}
	if oneTicketInfo.Type == "" {
		err = fmt.Errorf("消耗门票方式异常")
		return
	}
	if in.SubNum > oneTicketInfo.MaxBuyNum {
		err = fmt.Errorf("超过最大可认购数量")
		return
	}
	record, e := SubscribeRecord.GetSubscribeRecords(as.Id, in.UserId)
	if e != nil {
		err = e
		return
	}
	if record != nil {
		err = fmt.Errorf("已认购过")
		return
	}
	return
}

func (s *subscribeActivity) SetSubResult(in model.DoSubResult) {
	_, err := g.Redis().Do("SET", fmt.Sprintf(SubSetResultKey, in.OrderNo), gconv.String(in), "ex", 3600)
	if err != nil {
		g.Log().Errorf("DoSub err:%v", err)
		return
	}
	return
}

func (s *subscribeActivity) GetSubResult(orderNo string) (ret model.DoSubResult, err error) {
	gv, e := g.Redis().DoVar("GET", fmt.Sprintf(SubSetResultKey, orderNo))
	if e != nil {
		err = e
		g.Log().Errorf("DoSub err:%v", err)
		return
	}
	if gv == nil {
		err = fmt.Errorf("内部错误")
		return
	}
	err = gv.Scan(&ret)
	return
}

func (s *subscribeActivity) GetSubAwardResult(aid int, userId string) (ret *model.SubscribeRecord, err error) {
	err = g.DB().Model("subscribe_records").Where("aid = ? AND user_id = ? AND once_show = 0 AND award != 0", aid, userId).Scan(&ret)
	if err != nil {
		return
	}
	if ret != nil {
		_, err = g.DB().Exec("UPDATE subscribe_records SET once_show = 1 WHERE aid = ? AND user_id = ?", aid, userId)
		if err != nil {
			return
		}
	}
	return
}

func (s *subscribeActivity) DoSub(in model.DoSubReq) {
	orderNo := in.OrderNo
	//校验是否可认购
	oneTicketInfo, as, err := s.DoSubVerify(in)
	if err != nil {
		s.SetSubResult(model.DoSubResult{
			Reason:  err.Error(),
			Step:    "fail",
			Type:    in.Type,
			OrderNo: orderNo,
		})
		g.Log().Errorf("DoSub err:%v", err)
		return
	}
	if in.Type == model.TICKET_MONTH { //消耗月票方式认购
		tx, e := g.DB().Begin()
		if e != nil {
			s.SetSubResult(model.DoSubResult{
				Reason:  e.Error(),
				Step:    "fail",
				Type:    in.Type,
				OrderNo: orderNo,
			})
			g.Log().Errorf("DoSub err:%v", e)
			return
		}
		err = SubscribeRecord.CreateSubscribeRecord(tx, model.SubscribeRecord{
			ActivityType:       as.ActivityType,
			Aid:                as.Id,
			Name:               as.Name,
			Icon:               as.CoverImgUrl,
			UserId:             in.UserId,
			BuyNum:             in.SubNum,
			OrderNo:            orderNo,
			SumUnitMonthTicket: in.SubNum * oneTicketInfo.UnitNum,
			TicketType:         model.TICKET_MONTH,
			PayEndTime:         as.PayEndTime,
			PublisherId:        as.PublisherId,
		})
		if err != nil {
			tx.Rollback()
			s.SetSubResult(model.DoSubResult{
				Reason:  err.Error(),
				Step:    "fail",
				Type:    in.Type,
				OrderNo: orderNo,
			})
			g.Log().Errorf("DoSub err:%v", err)
			return
		}
		if oneTicketInfo.UnitNum > 0 {
			err = provider.User.OptUserMonthTicket(&map[string]interface{}{
				"userId": in.UserId,                         //用户ID
				"num":    oneTicketInfo.UnitNum * in.SubNum, //月票数量
				"type":   2,                                 //1.增加月票 2.减少月票
				"recordList": []map[string]interface{}{
					{
						"num":    oneTicketInfo.UnitNum * in.SubNum,                          //月票数量
						"source": 3,                                                          //来源 0.默认 1.持有资产 2.会员 3.优先购 4.投票
						"extra":  fmt.Sprintf(`{"desc":"元初发射台门票消耗","alias":"%s"}`, in.Alias), //扩展信息
					},
				},
			})
			if err != nil {
				tx.Rollback()
				s.SetSubResult(model.DoSubResult{
					Reason:  err.Error(),
					Step:    "fail",
					Type:    in.Type,
					OrderNo: orderNo,
				})
				g.Log().Errorf("DoSub err:%v", err)
				return
			}
		}
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			s.SetSubResult(model.DoSubResult{
				Reason:  err.Error(),
				Step:    "fail",
				Type:    in.Type,
				OrderNo: orderNo,
			})
			g.Log().Errorf("DoSub err:%v", err)
			return
		}
		s.SetSubResult(model.DoSubResult{
			Reason:  "success",
			Step:    "success",
			Type:    in.Type,
			OrderNo: orderNo,
		})
	} else if in.Type == model.TICKET_CRYSTAL { //消耗元晶方式
		tx, e := g.DB().Begin()
		if e != nil {
			s.SetSubResult(model.DoSubResult{
				Reason:  e.Error(),
				Step:    "fail",
				Type:    in.Type,
				OrderNo: orderNo,
			})
			g.Log().Errorf("DoSub err:%v", e)
			return
		}
		err = SubscribeRecord.CreateSubscribeRecord(tx, model.SubscribeRecord{
			ActivityType:   as.ActivityType,
			Aid:            as.Id,
			Name:           as.Name,
			Icon:           as.CoverImgUrl,
			UserId:         in.UserId,
			BuyNum:         in.SubNum,
			OrderNo:        orderNo,
			SumUnitCrystal: in.SubNum * oneTicketInfo.UnitNum,
			TicketType:     model.TICKET_CRYSTAL,
			PayEndTime:     as.PayEndTime,
			PublisherId:    as.PublisherId,
		})
		if err != nil {
			tx.Rollback()
			s.SetSubResult(model.DoSubResult{
				Reason:  err.Error(),
				Step:    "fail",
				Type:    in.Type,
				OrderNo: orderNo,
			})
			g.Log().Errorf("DoSub err:%v", err)
			return
		}
		if oneTicketInfo.UnitNum > 0 {
			err = provider.User.YJTransfer(&map[string]interface{}{
				"userId":   in.UserId,
				"category": 2,
				"amount":   oneTicketInfo.UnitNum * in.SubNum,
				"source":   23,
				"orderNo":  utils.Generate(),
			})
			if err != nil {
				tx.Rollback()
				s.SetSubResult(model.DoSubResult{
					Reason:  err.Error(),
					Step:    "fail",
					Type:    in.Type,
					OrderNo: orderNo,
				})
				g.Log().Errorf("DoSub err:%v", err)
				return
			}
		}
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			s.SetSubResult(model.DoSubResult{
				Reason:  err.Error(),
				Step:    "fail",
				Type:    in.Type,
				OrderNo: orderNo,
			})
			g.Log().Errorf("DoSub err:%v", err)
			return
		}
		s.SetSubResult(model.DoSubResult{
			Reason:  "success",
			Step:    "success",
			Type:    in.Type,
			OrderNo: orderNo,
		})
	} else if in.Type == model.TICKET_MONEY { //直接花钱
		//扩展参数
		extra := model.SubscribeRecordQueueData{}
		extra.Aid = as.Id
		extra.ActivityType = as.ActivityType
		extra.Name = as.Name
		extra.Icon = as.CoverImgUrl
		extra.UserId = in.UserId
		extra.BuyNum = in.SubNum
		extra.OrderNo = orderNo
		extra.SumUnitPrice = in.SubNum * oneTicketInfo.UnitNum
		extra.TicketType = model.TICKET_MONEY
		extra.PayEndTime = as.PayEndTime
		extra.FromUserId = in.UserId
		extra.ToUserId = "B"
		extra.PublisherId = as.PublisherId

		//向聚合支付下单
		orderReq := new(provider.CreateOrderReq)
		orderReq.ClientIp = in.ClientIp
		orderReq.UserId = in.UserId
		orderReq.AppType = "launchpad_ticket"
		orderReq.PayAmount = oneTicketInfo.UnitNum * in.SubNum
		orderReq.PayExpire = gtime.Now().Add(time.Minute * 10)
		orderReq.Subject = "购买元初发射台门票"
		orderReq.Description = "购买元初发射台门票"
		orderReq.SuccessRedirectUrl = in.SuccessRedirectUrl
		orderReq.ExitRedirectUrl = in.ExitRedirectUrl
		orderReq.AppOrderNo = orderNo
		orderReq.PublisherId = in.PublisherId
		orderReq.PlatformAppId = in.PlatformAppId
		orderReq.Extra = gconv.String(extra)
		err = provider.Payment.CreateOrder(orderReq)
		if err != nil {
			s.SetSubResult(model.DoSubResult{
				Reason:  "下单失败",
				Step:    "fail",
				Type:    in.Type,
				OrderNo: orderNo,
			})
			g.Log().Errorf("DoSub err:%v", err)
			return
		}
		s.SetSubResult(model.DoSubResult{
			Reason:  "success",
			Step:    "success",
			Type:    in.Type,
			OrderNo: orderNo,
		})
	}
	return
}

func (s *subscribeActivity) GetWaitOpenAwardActivity() (as []model.SubscribeActivity, err error) {
	err = g.DB().Model("subscribe_activity").Where("open_award_time <= ? and award_status = 0", time.Now()).Scan(&as)
	return
}

func (s *subscribeActivity) GetIngActivity() (as []model.SubscribeActivity, err error) {
	now := time.Now()
	g.DB().Model("subscribe_activity").Where("activity_start_time < ? AND ? < pay_end_time AND award_status = 2", now, now).Scan(&as)
	return
}

func (s *subscribeActivity) GetWaitPayEndActivity() (as []model.SubscribeActivity, err error) {
	err = g.DB().Model("subscribe_activity").Where("pay_end_time <= ? and award_status = 1 and pay_end = 0", time.Now()).Scan(&as)
	return
}

func (s *subscribeActivity) UpdatePayEnd(id int) (err error) {
	var ret sql.Result
	ret, err = g.DB().Exec("UPDATE subscribe_activity SET pay_end = 1 WHERE id = ?", id)
	if err != nil {
		return
	}
	affectedNum, _ := ret.RowsAffected()
	if affectedNum != 1 {
		err = fmt.Errorf("更新支付结束状态失败")
		return
	}
	return
}

func (s *subscribeActivity) UpdateActivityStatus(aid, status int) (err error) {
	var ret sql.Result
	if status == 2 {
		ret, err = g.DB().Exec("UPDATE subscribe_activity SET award_status = ?,award_complete_time = ? WHERE id = ?", status, time.Now(), aid)
		if err != nil {
			return
		}
	} else {
		ret, err = g.DB().Exec("UPDATE subscribe_activity SET award_status = ? WHERE id = ?", status, aid)
		if err != nil {
			return
		}
	}
	affectedNum, _ := ret.RowsAffected()
	if affectedNum != 1 {
		err = fmt.Errorf("更新活动抽奖状态失败")
		return
	}
	if status == model.AWARD_STATUS_END {
		go SubscribeRecord.SendSms(aid)
	}
	return
}

func (s *subscribeActivity) UpdateActivityRemainNum(tx *gdb.TX, aid, num int) (err error) {
	var ret sql.Result
	ret, err = tx.Exec("UPDATE subscribe_activity SET remain_num = remain_num - ? WHERE id = ?", num, aid)
	if err != nil {
		return
	}
	affectedNum, _ := ret.RowsAffected()
	if affectedNum != 1 {
		err = fmt.Errorf("更新失败")
		return
	}
	return
}

//更新活动为全部中签
func (s *subscribeActivity) AllAward(tx *gdb.TX, aid int) (err error) {
	_, err = tx.Exec("UPDATE subscribe_activity SET remain_num = remain_num - sub_sum WHERE id = ?", aid)
	return
}

//更新认购支付超时
func (s *subscribeActivity) UpdateSubPayTimeout(id int) (err error) {
	var result sql.Result
	result, err = g.DB().Exec("UPDATE subscribe_records SET pay_status = 2 WHERE id = ?", id)
	if err != nil {
		return
	}
	affectedNum, _ := result.RowsAffected()
	if affectedNum != 1 {
		err = fmt.Errorf("更新超时状态失败")
		return
	}
	return
}

//查询某人的超时次数
func (s *subscribeActivity) DoSubPayTimeOut() {
	var records []model.SubscribeRecord
	_ = g.DB().Model("subscribe_records").Where("award = 1 AND pay_end_time < ? AND pay_status = 0", time.Now()).Scan(&records)
	userMapNum := make(map[string]int)
	for _, v := range records {
		err := s.UpdateSubPayTimeout(v.Id)
		if err != nil {
			g.Log().Errorf("DoSubPayTimeOut err:%v", err)
			return
		}
		userMapNum[fmt.Sprintf(cache.SUB_PAY_TIMEOUT, v.UserId, v.ActivityType)]++
	}
	for k, num := range userMapNum {
		if num == 1 {
			g.Redis().Do("SET", k, num, "ex", 3600*24*30)
		} else {
			g.Redis().Do("SET", k, num, "ex", 3600*24*90)
		}
	}
	return
}

func (s *subscribeActivity) GetByIds(ids []int) (ret map[int]model.SubscribeActivity) {
	ret = make(map[int]model.SubscribeActivity)
	var as []model.SubscribeActivity
	g.DB().Model("subscribe_activity").Where("id IN (?)", ids).Scan(&as)
	for _, v := range as {
		ret[v.Id] = v
	}
	return
}

func (s *subscribeActivity) GetByAlias(alias string) (ret model.SubscribeActivity) {
	g.DB().Model("subscribe_activity").Where("alias = ?", alias).Scan(&ret)
	return
}
