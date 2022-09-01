package service

import (
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"meta_launchpad/model"
	"meta_launchpad/provider"
	"time"
)

const SubSetSeckillResultKey = "meta_launchpad:activity_seckill_result:%s"

type seckillActivity struct {
}

var SeckillActivity = new(seckillActivity)

func (s *seckillActivity) GetValidDetail(alias, publisherId string) (ret model.SeckillActivityFull, err error) {
	var as *model.SeckillActivity
	now := time.Now()
	err = g.DB().Model("seckill_activity").Where("alias = ?", alias).Scan(&as)
	if err != nil {
		return
	}
	if as == nil {
		err = fmt.Errorf("活动不存在")
		return
	}
	ret.SeckillActivity = as
	if now.Unix() > as.ActivityStartTime.Unix() && now.Unix() < as.ActivityEndTime.Unix() {
		ret.Status = model.SeckillActivityStatus_Ing
	} else {
		if now.Unix() < as.ActivityStartTime.Unix() {
			ret.Status = model.SeckillActivityStatus_Wait_Start
			ret.LastSec = as.ActivityStartTime.Unix() - now.Unix()
		} else {
			ret.Status = model.SeckillActivityStatus_End
		}
	}
	ret.PriceYuan = fmt.Sprintf("%.2f", float64(ret.Price)/100)

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

	ret.ChainName = publisherInfo.ChainName
	ret.ChainAddr = publisherInfo.ChainAddr
	ret.ChainType = publisherInfo.ChainType
	ret.AssetTotal = assetDetail.Total
	ret.AssetCreateAt = assetDetail.CreateTime
	ret.AssetDetailImg = templateInfo.DetailImg
	ret.NfrDay = as.NfrSec / 3600 / 24
	return
}

func (s *seckillActivity) SetSubResult(in model.DoSubResult) {
	_, err := g.Redis().Do("SET", fmt.Sprintf(SubSetSeckillResultKey, in.OrderNo), gconv.String(in), "ex", 3600)
	if err != nil {
		g.Log().Errorf("DoBuy err:%v", err)
		return
	}
	return
}

//回退库存
func (s *seckillActivity) UpdateRemain(tx *gdb.TX, aid int, num int) (err error) {
	_, err = tx.Exec("UPDATE seckill_activity SET remain_num = remain_num + ? WHERE id = ?", num, aid)
	return
}

func (s *seckillActivity) GetSubResult(orderNo string) (ret model.DoSubResult, err error) {
	gv, e := g.Redis().DoVar("GET", fmt.Sprintf(SubSetSeckillResultKey, orderNo))
	if e != nil {
		err = e
		g.Log().Errorf("DoBuy err:%v", err)
		return
	}
	if gv == nil {
		err = fmt.Errorf("内部错误")
		return
	}
	err = gv.Scan(&ret)
	return
}

func (s *seckillActivity) DoBuy(in model.DoBuyReq) {
	activityInfo, e := s.GetValidDetail(in.Alias, in.PublisherId)
	if e != nil {
		s.SetSubResult(model.DoSubResult{
			Reason:  e.Error(),
			Step:    "fail",
			OrderNo: in.OrderNo,
		})
		return
	}
	if activityInfo.Disable == 1 {
		s.SetSubResult(model.DoSubResult{
			Reason:  "活动已禁用",
			Step:    "fail",
			OrderNo: in.OrderNo,
		})
		return
	}
	if activityInfo.Status == model.SeckillActivityStatus_End {
		s.SetSubResult(model.DoSubResult{
			Reason:  "活动已结束",
			Step:    "fail",
			OrderNo: in.OrderNo,
		})
		return
	}
	//gv, e := g.Redis().DoVar("GET", fmt.Sprintf(cache.SECKILL_DISCIPLINE, in.UserId))
	//if e != nil {
	//	s.SetSubResult(model.DoSubResult{
	//		Reason:  e.Error(),
	//		Step:    "fail",
	//		OrderNo: in.OrderNo,
	//	})
	//	return
	//}
	//if !gv.IsEmpty() {
	//	s.SetSubResult(model.DoSubResult{
	//		Reason:  "您已超时一次未支付订单，暂不能参与秒杀活动",
	//		Step:    "fail",
	//		OrderNo: in.OrderNo,
	//	})
	//	return
	//}
	params := &map[string]interface{}{
		"appId":      activityInfo.AppId,
		"templateId": activityInfo.TemplateId,
	}
	assetInfo, e := provider.Asset.GetMateDataByAm(params)
	if e != nil {
		s.SetSubResult(model.DoSubResult{
			Reason:  e.Error(),
			Step:    "fail",
			OrderNo: in.OrderNo,
		})
		return
	}

	var tx *gdb.TX
	tx, e = g.DB().Begin()
	if e != nil {
		s.SetSubResult(model.DoSubResult{
			Reason:  e.Error(),
			Step:    "fail",
			OrderNo: in.OrderNo,
		})
		return
	}
	//扣除库存
	r, e := tx.Exec("UPDATE seckill_activity SET remain_num = remain_num - ? WHERE id = ?", in.Num, activityInfo.Id)
	if e != nil {
		tx.Rollback()
		s.SetSubResult(model.DoSubResult{
			Reason:  "库存不足",
			Step:    "fail",
			OrderNo: in.OrderNo,
		})
		return
	}
	affectedNum, _ := r.RowsAffected()
	if affectedNum != 1 {
		tx.Rollback()
		s.SetSubResult(model.DoSubResult{
			Reason:  "更新库存失败",
			Step:    "fail",
			OrderNo: in.OrderNo,
		})
		return
	}
	//刷新每个人可购买的数量
	e = SeckillUserBnum.CreateAndDecr(tx, model.SeckillUserBnum{
		Aid:    activityInfo.Id,
		UserId: in.UserId,
		CanBuy: activityInfo.LimitBuy,
	}, in.Num)
	if e != nil {
		tx.Rollback()
		s.SetSubResult(model.DoSubResult{
			Reason:  e.Error(),
			Step:    "fail",
			OrderNo: in.OrderNo,
		})
		return
	}
	//创建秒杀订单
	interOrder := model.SeckillOrder{
		OrderNo:     in.OrderNo,
		Num:         in.Num,
		RealFee:     in.Num * activityInfo.Price,
		UserId:      in.UserId,
		Aid:         activityInfo.Id,
		Name:        assetInfo.AssetName,
		Icon:        assetInfo.Icon,
		Status:      1,
		Price:       activityInfo.Price,
		PublisherId: in.PublisherId,
		PayExpireAt: gtime.Now().Add(time.Minute * 10),
	}
	e = SeckillOrder.Create(tx, interOrder)
	if e != nil {
		tx.Rollback()
		s.SetSubResult(model.DoSubResult{
			Reason:  e.Error(),
			Step:    "fail",
			OrderNo: in.OrderNo,
		})
		return
	}
	e = SeckillWaitPayOrder.Create(tx, model.SeckillWaitPayOrder{
		OrderNo:     interOrder.OrderNo,
		PayExpireAt: interOrder.PayExpireAt,
	})
	if e != nil {
		tx.Rollback()
		s.SetSubResult(model.DoSubResult{
			Reason:  e.Error(),
			Step:    "fail",
			OrderNo: in.OrderNo,
		})
		return
	}

	//扩展参数
	extra := model.SubscribeRecordQueueData{}
	extra.FromUserId = in.UserId
	extra.ToUserId = "B"
	extra.TotalFee = interOrder.RealFee

	orderReq := new(provider.CreateOrderReq)
	orderReq.ClientIp = in.ClientIp
	orderReq.UserId = in.UserId
	orderReq.AppType = "launchpad_seckill"
	orderReq.PayAmount = activityInfo.Price * in.Num
	orderReq.PayExpire = gtime.Now().Add(time.Minute * 10)
	orderReq.Subject = "秒杀活动"
	orderReq.Description = "秒杀活动"
	orderReq.SuccessRedirectUrl = in.SuccessRedirectUrl
	orderReq.ExitRedirectUrl = in.ExitRedirectUrl
	orderReq.AppOrderNo = interOrder.OrderNo
	orderReq.PublisherId = in.PublisherId
	orderReq.PlatformAppId = in.PlatformAppId
	orderReq.Extra = gconv.String(extra)
	e = provider.Payment.CreateOrder(orderReq)
	if e != nil {
		tx.Rollback()
		s.SetSubResult(model.DoSubResult{
			Reason:  "下单失败",
			Step:    "fail",
			OrderNo: in.OrderNo,
		})
		g.Log().Errorf("DoBuy err:%v", e)
		return
	}
	e = tx.Commit()
	if e != nil {
		tx.Rollback()
		s.SetSubResult(model.DoSubResult{
			Reason:  e.Error(),
			Step:    "fail",
			OrderNo: in.OrderNo,
		})
		g.Log().Errorf("DoBuy err:%v", e)
		return
	}
	s.SetSubResult(model.DoSubResult{
		Reason:  "success",
		Step:    "success",
		OrderNo: interOrder.OrderNo,
	})
	return
}

func (s *seckillActivity) GetSimpleDetail(aid int) (ret *model.SeckillActivity, err error) {
	err = g.DB().Model("seckill_activity").Where("id = ?", aid).Scan(&ret)
	if err != nil {
		return
	}
	if ret == nil {
		err = fmt.Errorf("活动不存在")
		return
	}
	return
}

func (s *seckillActivity) GetByIds(ids []int) (ret map[int]model.SeckillActivity) {
	ret = make(map[int]model.SeckillActivity)
	var as []model.SeckillActivity
	g.DB().Model("seckill_activity").Where("id IN (?)", ids).Scan(&as)
	for _, v := range as {
		ret[v.Id] = v
	}
	return
}
