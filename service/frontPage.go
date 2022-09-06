package service

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/model"
	"time"
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

func (c *frontPage) VolumeOfTrade(publisherId string, day int) (dealNum []model.Trade, payment []model.Trade) {
	t := time.Now()
	nowTime := t.AddDate(0, 0, 0).Format("2006-01-02 15:04:05")
	var sectionTime string
	var num int
	if day == 1 {
		sectionTime = t.AddDate(0, 0, -6).Format("2006-01-02")
		num = 7
	} else if day == 2 {
		sectionTime = t.AddDate(0, 0, -29).Format("2006-01-02")
		num = 30
	} else if day == 3 {
		sectionTime = t.AddDate(0, 0, -89).Format("2006-01-02")
		num = 90
	}
	// 成交笔数
	sql := "SELECT t0.date,IFNULL(t1.count,0) count FROM (SELECT @cdate := DATE_ADD(@cdate, INTERVAL + 1 DAY) date FROM (SELECT @cdate := DATE_ADD('" + sectionTime + "', INTERVAL - 1 DAY) date FROM subscribe_records) l) t0 LEFT JOIN (SELECT DATE_ADD(DATE_FORMAT(created_at,'%Y-%m-%d'), INTERVAL 0 DAY) created_at ,COUNT(1) count "
	sql += fmt.Sprintf("FROM subscribe_records WHERE pay_status = 1 AND created_at BETWEEN '%s' AND '%s' AND publisher_id = '%s'", sectionTime, nowTime, publisherId)
	sql += " GROUP BY DATE_FORMAT(created_at,'%Y-%m-%d') UNION SELECT DATE_FORMAT(created_at,'%Y-%m-%d') created_at ,COUNT(1) count FROM seckill_orders "
	sql += fmt.Sprintf("WHERE `status` = 1 AND created_at BETWEEN '%s' AND '%s' AND publisher_id = '%s' ", sectionTime, nowTime, publisherId)
	sql += " GROUP BY DATE_FORMAT(created_at,'%Y-%m-%d')) t1 on t0.date = t1.created_at ORDER BY t0.date"
	sql += fmt.Sprintf(" LIMIT %d", num)
	err := g.DB().GetScan(&dealNum, sql)
	if err != nil {
		return
	}

	// 支付人数
	paySql := "SELECT t0.date,IFNULL(t1.count,0) count FROM (SELECT @cdate := DATE_ADD(@cdate, INTERVAL + 1 DAY) date FROM (SELECT @cdate := DATE_ADD('" + sectionTime + "', INTERVAL - 1 DAY) date FROM subscribe_records) l) t0 LEFT JOIN ( SELECT DATE_FORMAT(created_at,'%Y-%m-%d') created_at,COUNT(1) count"
	paySql += fmt.Sprintf(" from (SELECT user_id,created_at FROM subscribe_records WHERE pay_status = 1 AND created_at BETWEEN '%s' AND '%s' AND publisher_id = '%s' GROUP BY user_id UNION SELECT user_id,created_at FROM seckill_orders WHERE `status` = 1 AND created_at BETWEEN '%s' AND '%s' AND publisher_id = '%s' GROUP BY user_id) l", sectionTime, nowTime, publisherId, sectionTime, nowTime, publisherId)
	paySql += " GROUP BY DATE_FORMAT(created_at,'%Y-%m-%d')"
	paySql += fmt.Sprintf(") t1 on t0.date = t1.created_at ORDER BY t0.date LIMIT %d", num)
	err = g.DB().GetScan(&payment, paySql)
	if err != nil {
		return
	}
	return
}
