package service

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	"database/sql"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"meta_launchpad/cache"
	"meta_launchpad/model"
	"meta_launchpad/provider"
	"time"
)

type subscribeRecord struct {
}

var SubscribeRecord = new(subscribeRecord)

func (s *subscribeRecord) GetSubscribeAwardRecord(alias string) (ret model.GetSubscribeAwardRecordRet, err error) {
	as, e := SubscribeActivity.GetSimpleDetailByAlias(alias)
	if e != nil {
		err = e
		return
	}
	if as == nil {
		err = fmt.Errorf("活动不存在")
		return
	}
	if as.AwardStatus != 2 {
		return
	}
	if as.PayEndTime.Unix() < time.Now().Unix() {
		return
	}
	var records []model.SubscribeRecord
	err = g.DB().Model("subscribe_records").Where("aid = ? AND award_num != 0", as.Id).Scan(&records)
	if err != nil {
		return
	}
	ret.AwardCompleteTime = as.AwardCompleteTime
	userIds := make([]string, 0)
	for _, v := range records {
		userIds = append(userIds, v.UserId)
	}
	_, uMap, _ := provider.User.GetUserInfo(userIds)
	for _, v := range records {
		ret.List = append(ret.List, model.GetSubscribeAwardRecordRetItem{
			Phone:    utils.FormatMobileStar(uMap[v.UserId].Phone),
			AwardNum: v.AwardNum,
		})
	}
	return
}

func (s *subscribeRecord) GetSubscribeRecords(aid int, userId string) (ret *model.SubscribeRecord, err error) {
	err = g.DB().Model("subscribe_records").Where("aid = ? AND user_id = ?", aid, userId).Scan(&ret)
	if err != nil {
		return
	}
	return
}

func (s *subscribeRecord) CreateSubscribeRecord(tx *gdb.TX, in model.SubscribeRecord) (err error) {
	_, err = tx.Model("subscribe_records").Insert(&in)
	if err != nil {
		return
	}
	var result sql.Result
	result, err = tx.Exec("UPDATE subscribe_activity SET sub_sum_people = sub_sum_people + 1,sub_sum = sub_sum + ? WHERE id = ?", in.BuyNum, in.Aid)
	if err != nil {
		return
	}
	affectedNum, _ := result.RowsAffected()
	if affectedNum != 1 {
		err = fmt.Errorf("操作失败")
		return
	}
	return
}

//获取未公布中签结果的用户
func (s *subscribeRecord) GetWaitLuckyDraw(aid int) (ret []model.SubscribeRecord, err error) {
	err = g.DB().Model("subscribe_records").Where("aid = ? and award = 0", aid).Scan(&ret)
	if err != nil {
		return
	}
	return
}

//更新活动为全部中签
func (s *subscribeRecord) AllAward(tx *gdb.TX, aid int, unitPrice int) (err error) {
	_, err = tx.Exec("UPDATE subscribe_records SET award = 1,award_num = buy_num,award_at = ?,sum_price = buy_num*? WHERE aid = ? and award = 0", gtime.Now(), unitPrice, aid)
	return
}

//更新活动为全部未中签
func (s *subscribeRecord) AllUnAward(aid int) (err error) {
	_, err = g.DB().Exec("UPDATE subscribe_records SET award = 2,award_at = ? WHERE aid = ? and award = 0", gtime.Now(), aid)
	return
}

//更新中签
func (s *subscribeRecord) UpdateAward(tx *gdb.TX, id, awardNum int, unitPrice int) (err error) {
	_, err = tx.Exec("UPDATE subscribe_records SET award = 1,award_num = ?,award_at = ?,sum_price = ? WHERE id = ?", awardNum, gtime.Now(), awardNum*unitPrice, id)
	return
}

//更新未中签
func (s *subscribeRecord) UpdateUnAward(id int) (err error) {
	_, err = g.DB().Exec("UPDATE subscribe_records SET award = 2,award_at = ? WHERE aid = ? AND award_num = 0", gtime.Now(), id)
	return
}

//获取到未全部中签的人（包括了未中签的）
func (s *subscribeRecord) GetUnFullAward(aid int) (ret []model.SubscribeRecord, err error) {
	err = g.DB().Model("subscribe_records").Where("aid = ? and award_num < buy_num", aid).Scan(&ret)
	return
}

//获取认购记录列表
func (s *subscribeRecord) GetList(userId string, publisherId string, pageNum int, award int) (ret model.SubscribeRecordList, err error) {
	pageSize := 20
	m := g.DB().Model("subscribe_records").Where("user_id = ? and publisher_id = ?", userId, publisherId)
	if award != -1 {
		m = m.Where("award", award)
	}
	ret.Total, err = m.Count()
	if err != nil {
		return
	}
	if ret.Total == 0 {
		return
	}
	m.Page(pageNum, pageSize).Order("id desc").Scan(&ret.List)
	return
}

//获取认购记录详情
func (s *subscribeRecord) GetDetail(orderNo string) (ret model.SubscribeRecordDetail, err error) {
	var record *model.SubscribeRecord
	err = g.DB().Model("subscribe_records").Where("order_no = ?", orderNo).Scan(&record)
	if err != nil {
		return
	}
	if record == nil {
		err = fmt.Errorf("未查询到信息")
		return
	}
	if record.TicketType == model.TICKET_MONTH {
		ret.ConsumeUnit = fmt.Sprintf("%d月票", record.SumUnitMonthTicket)
	}
	if record.TicketType == model.TICKET_CRYSTAL {
		ret.ConsumeUnit = fmt.Sprintf("%d元晶", record.SumUnitCrystal)
	}
	if record.TicketType == model.TICKET_MONEY {
		ret.ConsumeUnit = fmt.Sprintf("%.2f元", float64(record.SumUnitPrice)/100)
	}
	ret.AwardNum = record.AwardNum
	ret.Name = record.Name
	ret.AwardAt = record.AwardAt
	ret.Award = record.Award
	ret.BuyNum = record.BuyNum
	ret.Icon = record.Icon
	ret.UserId = record.UserId
	ret.Aid = record.Aid
	ret.CreatedAt = record.CreatedAt
	return
}

//获取认购记录详情
func (s *subscribeRecord) GetSimpleDetail(orderNo string) (ret *model.SubscribeRecord, err error) {
	err = g.DB().Model("subscribe_records").Where("order_no = ?", orderNo).Scan(&ret)
	if err != nil {
		return
	}
	if ret == nil {
		err = fmt.Errorf("未查询到信息")
		return
	}
	return
}

//以订单维度查询认购记录
func (s *subscribeRecord) GetListByOrder(userId string, orderNo string, pageNum int, status int, publisherId string, activityType int) (ret model.SubscribeListByOrderRet, err error) {
	var records []model.SubscribeRecord
	m := g.DB().Model("subscribe_records").Where("user_id = ? AND award = 1 AND publisher_id = ?", userId, publisherId)
	if activityType != 0 {
		m = m.Where("activity_type = ?", activityType)
	}
	if status != -1 {
		m = m.Where("pay_status = ?", status)
	}
	if orderNo != "" {
		m = m.Where("(order_no = ? or pay_order_no = ?)", orderNo, orderNo)
	}
	ret.Total, err = m.Count()
	if err != nil {
		return
	}
	if ret.Total == 0 {
		return
	}
	pageSize := 20
	err = m.Page(pageNum, pageSize).Order("id desc").Scan(&records)
	if err != nil {
		return
	}
	for _, v := range records {
		lastSec := v.PayEndTime.Unix() - time.Now().Unix()
		if lastSec <= 0 {
			lastSec = 0
		}
		item := model.SubscribeListByOrderRetItem{
			BuyNum:        v.AwardNum,
			UnitPriceYuan: fmt.Sprintf("%.2f", float64(v.SumPrice)/float64(v.AwardNum)/100),
			SumPriceYuan:  fmt.Sprintf("%.2f", float64(v.SumPrice)/100),
			SumPrice:      v.SumPrice,
			OrderNo:       v.OrderNo,
			Name:          v.Name,
			Icon:          v.Icon,
			PayOrderNo:    v.PayOrderNo,
			PayEndTime:    v.PayEndTime,
			PaidAt:        v.PaidAt,
			PayMethod:     v.PayMethod,
			Status:        v.PayStatus,
			LastSec:       lastSec,
		}
		ret.List = append(ret.List, item)
	}
	return
}

//创建支付订单
func (s *subscribeRecord) CreateOrder(userId, clientIp, orderNo, successRedirectUrl, exitRedirectUrl, publisherId, appId string) (orderReq *provider.CreateOrderReq, err error) {
	info, e := s.GetListByOrder(userId, orderNo, 1, 0, publisherId, 0)
	if e != nil {
		err = e
		return
	}
	if len(info.List) == 0 {
		err = fmt.Errorf("订单不存在")
		return
	}
	record := info.List[0]
	//向聚合支付下单
	orderReq = new(provider.CreateOrderReq)
	orderReq.ClientIp = clientIp
	orderReq.UserId = userId
	orderReq.AppType = "launchpad_pay"
	orderReq.PayAmount = record.SumPrice
	orderReq.PayExpire = gtime.Now().Add(time.Minute * 10)
	orderReq.Subject = "元初发射台付款"
	orderReq.Description = "元初发射台付款"
	orderReq.SuccessRedirectUrl = successRedirectUrl
	orderReq.ExitRedirectUrl = exitRedirectUrl
	orderReq.PublisherId = publisherId
	orderReq.PlatformAppId = appId
	orderReq.AppOrderNo = fmt.Sprintf("%d", utils.GetOrderNo())
	orderReq.Extra = fmt.Sprintf(`{"fromUserId":"%s","toUserId":"B","orderNo":"%s","totalFee":%d}`, userId, orderNo, record.SumPrice)
	err = provider.Payment.CreateOrder(orderReq)
	if err != nil {
		g.Log().Errorf("CreateOrder err:%v", err)
		return
	}
	return
}

//已发放资产
func (s *subscribeRecord) UpdatePublishAsset(orderNo string) (err error) {
	var retSql sql.Result
	retSql, err = g.DB().Exec("UPDATE subscribe_records SET publish_asset = 1 WHERE order_no = ?", orderNo)
	if err != nil {
		return
	}
	affectedNum, _ := retSql.RowsAffected()
	if affectedNum != 1 {
		err = fmt.Errorf("更新状态失败")
		return
	}
	return
}

//已支付
func (s *subscribeRecord) DoPaid(payMethod string, orderNo string, payOrderNo string) (err error) {
	var retSql sql.Result
	retSql, err = g.DB().Exec("UPDATE subscribe_records SET pay_status = 1,paid_at = ?,pay_method = ?,pay_order_no = ? WHERE order_no = ?", gtime.Now(), payMethod, payOrderNo, orderNo)
	if err != nil {
		return
	}
	affectedNum, _ := retSql.RowsAffected()
	if affectedNum != 1 {
		err = fmt.Errorf("更新状态失败")
		return
	}
	return
}

//查出所有根据活动id，查出所有中签人，发送短信
func (s *subscribeRecord) SendSms(aid int) {
	var records []model.SubscribeRecord
	g.DB().Model("subscribe_records").Where("aid = ? AND award = 1", aid).Scan(&records)
	for _, r := range records {
		_, userMap, err := provider.User.GetUserInfo([]string{r.UserId})
		if err != nil {
			continue
		}
		_ = Sms.SendSms(userMap[r.UserId].Phone, SmsConfig[r.PublisherId], "a1609CKE", "4HZdAzLt", map[string]string{
			"googs": r.Name,
			"time":  r.PayEndTime.Layout("2006-01-02 15:04:05"),
		})
	}
}

//支付剩余15分钟，发送短信提醒
func (s *subscribeRecord) SendSmsWaitPay(as model.SubscribeActivity) {
	lock := cache.DistributedLock("SendSmsWaitPay_" + as.Alias)
	if lock {
		now := time.Now()
		if as.PayEndTime.Unix()-now.Unix() < 15*60 && (as.PayEndTime.Unix()-now.Unix() > 0) { //付款截止小于15分钟
			var records []model.SubscribeRecord
			g.DB().Model("subscribe_records").Where("aid = ? AND award = 1 AND pay_status = 0", as.Id).Scan(&records)
			for _, r := range records {
				_, userMap, err := provider.User.GetUserInfo([]string{r.UserId})
				if err != nil {
					continue
				}
				_ = Sms.SendSms(userMap[r.UserId].Phone, SmsConfig[r.PublisherId], "aIIbedlG", "4HZdAzLt", map[string]string{
					"goods": r.Name,
					"time":  r.PayEndTime.Layout("15:04"),
				})
			}
		}
	}
}
