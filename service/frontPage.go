package service

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/model"
	"strconv"
	"time"
)

type frontPage struct {
}

var FrontPage = new(frontPage)

func (c *frontPage) TransactionSlip(publisherId string) (transactionSlip []model.TransactionSlip, sum model.TransactionSlip) {
	sql := fmt.Sprintf("SELECT '优先购' `name`,COUNT(1) value FROM subscribe_records WHERE pay_status = 1 AND activity_type = 1 AND publisher_id = '%s' UNION SELECT '普通购' `name`,COUNT(1) count FROM subscribe_records WHERE pay_status = 1 AND activity_type = 2 AND publisher_id = '%s' UNION SELECT '秒杀购', COUNT(1) FROM seckill_orders WHERE `status` = 2 AND publisher_id = '%s' ", publisherId, publisherId, publisherId)
	db := g.DB()
	err := db.GetScan(&transactionSlip, sql)
	if err != nil {
		return
	}
	var count int
	for _, i := range transactionSlip {
		count += i.Value
	}
	sum.Name = "总单数"
	sum.Value = count
	return
}

// VolumeOfTrade 支付数，人数
func (c *frontPage) VolumeOfTrade(publisherId string, day int) (dealTime, paymentTime []string, dealCount, paymentCount []int) {
	var dealNum []model.Trade
	var payment []model.Trade
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
	sql := "SELECT t0.date created_at,IFNULL(t1.count,0) count FROM (SELECT @cdate := DATE_ADD(@cdate, INTERVAL + 1 DAY) date FROM (SELECT @cdate := DATE_ADD('" + sectionTime + "', INTERVAL - 1 DAY) date FROM subscribe_records) l) t0 LEFT JOIN (SELECT DATE_ADD(DATE_FORMAT(created_at,'%Y-%m-%d'), INTERVAL 0 DAY) created_at ,COUNT(1) count "
	sql += fmt.Sprintf("FROM subscribe_records WHERE pay_status = 1 AND created_at BETWEEN '%s' AND '%s' AND publisher_id = '%s'", sectionTime, nowTime, publisherId)
	sql += " GROUP BY DATE_FORMAT(created_at,'%Y-%m-%d') UNION SELECT DATE_FORMAT(created_at,'%Y-%m-%d') created_at ,COUNT(1) count FROM seckill_orders "
	sql += fmt.Sprintf("WHERE `status` = 2 AND created_at BETWEEN '%s' AND '%s' AND publisher_id = '%s' ", sectionTime, nowTime, publisherId)
	sql += " GROUP BY DATE_FORMAT(created_at,'%Y-%m-%d')) t1 on t0.date = t1.created_at ORDER BY t0.date"
	sql += fmt.Sprintf(" LIMIT %d", num)
	err := g.DB().GetScan(&dealNum, sql)
	if err != nil {
		return
	}
	for _, i := range dealNum {
		dealTime = append(dealTime, i.CreatedAt)
		dealCount = append(dealCount, i.Count)
	}

	// 支付人数
	paySql := "SELECT t0.date created_at,IFNULL(t1.count,0) count FROM (SELECT @cdate := DATE_ADD(@cdate, INTERVAL + 1 DAY) date FROM (SELECT @cdate := DATE_ADD('" + sectionTime + "', INTERVAL - 1 DAY) date FROM subscribe_records) l) t0 LEFT JOIN ( SELECT DATE_FORMAT(created_at,'%Y-%m-%d') created_at,COUNT(1) count"
	paySql += fmt.Sprintf(" from (SELECT user_id,created_at FROM subscribe_records WHERE pay_status = 1 AND created_at BETWEEN '%s' AND '%s' AND publisher_id = '%s' GROUP BY user_id UNION SELECT user_id,created_at FROM seckill_orders WHERE `status` = 2 AND created_at BETWEEN '%s' AND '%s' AND publisher_id = '%s' GROUP BY user_id) l", sectionTime, nowTime, publisherId, sectionTime, nowTime, publisherId)
	paySql += " GROUP BY DATE_FORMAT(created_at,'%Y-%m-%d')"
	paySql += fmt.Sprintf(") t1 on t0.date = t1.created_at ORDER BY t0.date LIMIT %d", num)
	err = g.DB().GetScan(&payment, paySql)
	for _, i := range payment {
		paymentTime = append(paymentTime, i.CreatedAt)
		paymentCount = append(paymentCount, i.Count)
	}
	if err != nil {
		return
	}
	return
}

// TransactionsNum 支付笔数
func (c *frontPage) TransactionsNum(publisherId string) (count int, float float64, err error) {
	sql := fmt.Sprintf("SELECT count(1) FROM (SELECT user_id,created_at FROM subscribe_records WHERE pay_status = 1 AND publisher_id = '%s' UNION SELECT user_id,created_at FROM seckill_orders WHERE `status` = 2 AND publisher_id = '%s' ) t", publisherId, publisherId)
	rows, err := g.DB().Query(sql)
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			g.Log().Error(err)
			return
		}
	}
	float = c.percentage("SELECT count(1) count FROM (SELECT COUNT(1) count FROM subscribe_records WHERE pay_status = 1 AND created_at BETWEEN '%s' AND '%s' AND publisher_id = '%s' UNION SELECT COUNT(1) count FROM seckill_orders WHERE `status` = 2 AND created_at BETWEEN '%s' AND '%s' AND publisher_id = '%s') t", publisherId)
	return
}

// Payers 支付人数
func (c *frontPage) Payers(publisherId string) (count int, float float64, err error) {
	sql := fmt.Sprintf("SELECT COUNT(1) count from (SELECT user_id,created_at FROM subscribe_records WHERE pay_status = 1 AND publisher_id = '%s' GROUP BY user_id UNION SELECT user_id,created_at FROM seckill_orders WHERE `status` = 1 AND publisher_id = '%s' GROUP BY user_id) l", publisherId, publisherId)
	rows, err := g.DB().Query(sql)
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			g.Log().Error(err)
			return
		}
	}
	float = c.percentage("SELECT COUNT(1) count from (SELECT user_id,created_at FROM subscribe_records WHERE pay_status = 1 AND created_at BETWEEN '%s' AND '%s' AND publisher_id = '%s' GROUP BY user_id UNION SELECT user_id,created_at FROM seckill_orders WHERE `status` = 1 AND created_at BETWEEN '%s' AND '%s' AND publisher_id = '%s' GROUP BY user_id) t", publisherId)
	return
}

// Turnover 成交额
func (c *frontPage) Turnover(publisherId string) (priceFloat []float64, priceTime []string, float float64, count float64, err error) {
	var price []model.Price
	t := time.Now()
	nowTime := t.AddDate(0, 0, 0).Format("2006-01-02 15:04:05")
	sectionTime := t.AddDate(0, 0, -6).Format("2006-01-02")
	sql := fmt.Sprintf("SELECT t0.date created_at,t1.price/100 price FROM (SELECT @cdate := DATE_ADD(@cdate, INTERVAL + 1 DAY) date FROM (SELECT @cdate := DATE_ADD('%s', INTERVAL - 1 DAY) date FROM subscribe_records) l) t0 LEFT JOIN (", sectionTime)
	sql += "SELECT DATE_FORMAT(created_at,'%Y-%m-%d') created_at,sum(sum_price) price FROM subscribe_records WHERE pay_status = 1 "
	sql += fmt.Sprintf("AND created_at BETWEEN '%s' AND '%s' AND publisher_id = '%s' ", sectionTime, nowTime, publisherId)
	sql += "GROUP BY DATE_FORMAT(created_at,'%Y-%m-%d') UNION SELECT DATE_FORMAT(created_at,'%Y-%m-%d'),sum(price) price FROM seckill_orders WHERE `status` = 2 "
	sql += fmt.Sprintf("AND created_at BETWEEN '%s' AND '%s' AND publisher_id = '%s' ", sectionTime, nowTime, publisherId)
	sql += "GROUP BY DATE_FORMAT(created_at,'%Y-%m-%d')) t1 ON t0.date=t1.created_at ORDER BY t0.date LIMIT 7"
	err = g.DB().GetScan(&price, sql)
	if err != nil {
		return
	}
	for _, i := range price {
		priceFloat = append(priceFloat, i.Price)
		priceTime = append(priceTime, i.CreatedAt)
	}

	float = c.percentage("SELECT sum(price)/100 count FROM (SELECT sum(sum_price) price,count(1) count FROM subscribe_records WHERE pay_status = 1 AND created_at BETWEEN '%s' AND '%s' AND publisher_id = '%s' UNION SELECT sum(price) price,count(1) FROM seckill_orders WHERE `status` = 2 AND created_at BETWEEN '%s' AND '%s' AND publisher_id = '%s' ) l", publisherId)

	countSql := fmt.Sprintf("SELECT sum(price)/100 FROM (SELECT sum(sum_price) price,count(1) count FROM subscribe_records WHERE pay_status = 1 AND publisher_id = '%s' UNION SELECT sum(price) price,count(1) FROM seckill_orders WHERE `status` = 2 AND publisher_id = '%s' ) l", publisherId, publisherId)
	query, err := g.DB().Query(countSql)
	if err != nil {
		return
	}
	for query.Next() {
		err = query.Scan(&count)
		if err != nil {
			return
		}
	}
	return
}

// percentage 对比上个月比例
func (c *frontPage) percentage(sql, publisherId string) (float float64) {
	type s struct {
		Count float64 `json:"count"`
	}
	t := time.Now()
	var percentage s
	oneMonth := t.AddDate(0, 0, -30).Format("2006-01-02") + " 00:00:00"
	nowDay := t.AddDate(0, 0, 0).Format("2006-01-02 ") + " 23:59:59"
	percentageSql := fmt.Sprintf(sql, oneMonth, nowDay, publisherId, oneMonth, nowDay, publisherId)
	err := g.DB().GetScan(&percentage, percentageSql)
	if err != nil {
		return
	}

	var percentage2 s
	towMonth := t.AddDate(0, 0, -60).Format("2006-01-02") + " 00:00:00"
	oneMonths := t.AddDate(0, 0, -31).Format("2006-01-02 ") + " 23:59:59"
	percentageSql2 := fmt.Sprintf(sql, towMonth, oneMonths, publisherId, towMonth, oneMonths, publisherId)
	err = g.DB().GetScan(&percentage2, percentageSql2)
	if err != nil {
		return
	}
	addPercentage := (percentage.Count - percentage2.Count) / percentage2.Count * 100
	float, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", addPercentage), 1)
	return
}
