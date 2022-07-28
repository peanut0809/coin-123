package task

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/cache"
	"meta_launchpad/model"
	"meta_launchpad/service"
	"time"
)

func LuckyDraw() {
	as, err := service.SubscribeActivity.GetWaitOpenAwardActivity()
	if err != nil {
		g.Log().Errorf("RunLuckyDrawTask err:%v", err)
		return
	}
	for _, v := range as {
		err = service.SubscribeActivity.UpdateActivityStatus(v.Id, model.AWARD_STATUS_ING)
		if err != nil {
			g.Log().Errorf("RunLuckyDrawTask err:%v", err)
			return
		}
		//库存为0，直接全部不中奖
		if v.RemainNum == 0 {
			err = service.SubscribeRecord.AllUnAward(v.Id)
			if err != nil {
				g.Log().Errorf("RunLuckyDrawTask err:%v", err)
				return
			}
			err = service.SubscribeActivity.UpdateActivityStatus(v.Id, model.AWARD_STATUS_END)
			if err != nil {
				g.Log().Errorf("RunLuckyDrawTask err:%v", err)
				return
			}
		} else {
			records, err := service.SubscribeRecord.GetWaitLuckyDraw(v.Id)
			if err != nil {
				g.Log().Errorf("RunLuckyDrawTask err:%v", err)
				return
			}
			if v.SubSum <= v.SumNum { //认购总数小于或等于发行总数，全部中签
				tx, e := g.DB().Begin()
				if e != nil {
					g.Log().Errorf("RunLuckyDrawTask err:%v", e)
					return
				}
				err = service.SubscribeActivity.AllAward(tx, v.Id)
				if err != nil {
					tx.Rollback()
					g.Log().Errorf("RunLuckyDrawTask err:%v", err)
					return
				}
				err = service.SubscribeRecord.AllAward(tx, v.Id, v.Price)
				if err != nil {
					tx.Rollback()
					g.Log().Errorf("RunLuckyDrawTask err:%v", err)
					return
				}
				err = tx.Commit()
				if err != nil {
					tx.Rollback()
					g.Log().Errorf("RunLuckyDrawTask err:%v", err)
					return
				}
				err = service.SubscribeActivity.UpdateActivityStatus(v.Id, model.AWARD_STATUS_END)
				if err != nil {
					g.Log().Errorf("RunLuckyDrawTask err:%v", err)
					return
				}
			} else {
				if v.ActivityType == 1 { //优先购
					if v.AwardMethod == 0 { //比例分配方式，库存一定够分
						err = LuckyDrawBySplit(v, records)
						if err != nil {
							g.Log().Errorf("RunLuckyDrawTask err:%v", err)
							return
						}
					}
					if v.AwardMethod == 1 { //抽签
						LuckyDrawByRandom(records, v.RemainNum, v.Id, v.Price)
					}
				}
				if v.ActivityType == 2 { //普通购
					LuckyDrawByRandom(records, v.RemainNum, v.Id, v.Price)
				}
				err = service.SubscribeActivity.UpdateActivityStatus(v.Id, model.AWARD_STATUS_END)
				if err != nil {
					g.Log().Errorf("RunLuckyDrawTask err:%v", err)
					return
				}
			}
		}
	}
}

//中签
func Award(rid, aid int, awardNum, unitPrice int) (err error) {
	tx, e := g.DB().Begin()
	if e != nil {
		err = e
		return
	}
	err = service.SubscribeActivity.UpdateActivityRemainNum(tx, aid, awardNum)
	if err != nil {
		tx.Rollback()
		return
	}
	err = service.SubscribeRecord.UpdateAward(tx, rid, awardNum, unitPrice)
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
	}
	return
}

//比例分配方式
func LuckyDrawBySplit(as model.SubscribeActivity, records []model.SubscribeRecord) (err error) {
	sumAwardNum := 0
	for _, v := range records {
		awardNum := v.BuyNum / as.SubSum * as.SumNum
		if awardNum != 0 {
			err = Award(v.Id, as.Id, awardNum, as.Price)
			if err == nil {
				sumAwardNum += awardNum
			} else {
				g.Log().Errorf("RunLuckyDrawTask err:%v", err)
				return
			}
		}
	}
	remainNum := as.SumNum - sumAwardNum
	if remainNum > 0 { //剩下的按签抽
		rs, e := service.SubscribeRecord.GetUnFullAward(as.Id)
		if e != nil {
			g.Log().Errorf("RunLuckyDrawTask err:%v", e)
			return
		}
		LuckyDrawByRandom(rs, remainNum, as.Id, as.Price)
	}
	return
}

func LuckyDrawByRandom(records []model.SubscribeRecord, remainNum, aid int, unitPrice int) {
	//取出
	var fullIndexs []int //总签数
	for i, v := range records {
		aid = v.Aid
		for j := 0; j < v.BuyNum-v.AwardNum; j++ {
			fullIndexs = append(fullIndexs, i)
		}
	}
	//洗牌
	fullIndexs = utils.Shuffle(fullIndexs)
	recordMap := make(map[int]int) //每个记录的中签数
	for _, i := range fullIndexs {
		if remainNum <= 0 {
			break
		}
		remainNum--
		recordMap[i]++
	}
	for i, awardNum := range recordMap {
		record := records[i]
		err := Award(record.Id, record.Aid, record.AwardNum+awardNum, unitPrice)
		if err != nil {
			g.Log().Errorf("RunLuckyDrawTask err:%v", err)
			return
		}
	}
	//更新未中奖记录
	err := service.SubscribeRecord.UpdateUnAward(aid)
	if err != nil {
		g.Log().Errorf("RunLuckyDrawTask err:%v", err)
		return
	}
}

const TASK_RunLuckyDrawTask = "RunLuckyDrawTask"
const TASK_CheckSubPayTimeout = "CheckSubPayTimeout"

//检查超时未支付
func CheckSubPayTimeout() {
	cache.DistributedUnLock(TASK_CheckSubPayTimeout)
	for {
		lock := cache.DistributedLock(TASK_CheckSubPayTimeout)
		if lock {
			service.SubscribeActivity.DoSubPayTimeOut()
			cache.DistributedUnLock(TASK_CheckSubPayTimeout)
		}
		time.Sleep(time.Second * 10)
	}
}

func RunLuckyDrawTask() {
	cache.DistributedUnLock(TASK_RunLuckyDrawTask)
	for {
		lock := cache.DistributedLock(TASK_RunLuckyDrawTask)
		if lock {
			LuckyDraw()
			cache.DistributedUnLock(TASK_RunLuckyDrawTask)
		}
		time.Sleep(time.Second * 10)
	}
}
