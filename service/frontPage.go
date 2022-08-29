package service

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/model"
)

type frontPage struct {
}

var FrontPage = new(frontPage)

func (c *frontPage) TransactionSlip(publisherId string) (transactionSlip []model.TransactionSlip, sum model.TransactionSlip) {
	sql := fmt.Sprintf("SELECT '优先购' `name`,COUNT(1) count FROM subscribe_records WHERE pay_status = 1 AND activity_type = 1 AND publisher_id = '%s' UNION SELECT '普通购' `name`,COUNT(1) count FROM subscribe_records WHERE pay_status = 1 AND activity_type = 2 AND publisher_id = '%s' UNION SELECT '秒杀购', COUNT(1) FROM seckill_orders WHERE `status` = 2 AND publisher_id = '%s' ", publisherId, publisherId, publisherId)
	db := g.DB()
	err := db.GetScan(&transactionSlip, sql)
	if err != nil {
		return
	}
	var count int
	for _, i := range transactionSlip {
		count += i.Count
	}
	sum.Name = "总单数"
	sum.Count = count
	return
}

func (c *frontPage) VolumeOfTrade(publisherId string) (count float64) {
	type s struct {
		VolumeOfTrade float64 `json:"volumeOfTrade"`
	}
	var volumeOfTrade []s
	sql := fmt.Sprintf("SELECT sum_price/100 volumeOfTrade FROM subscribe_records WHERE pay_status = 1 AND publisher_id = '%s' UNION SELECT real_fee/100 FROM seckill_orders WHERE `status` = 2 AND publisher_id = '%s' ", publisherId, publisherId)
	db := g.DB()
	err := db.GetScan(&volumeOfTrade, sql)
	if err != nil {
		return
	}
	for _, i := range volumeOfTrade {
		count += i.VolumeOfTrade
	}
	return
}
