package service

import (
	"fmt"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/model"
	"meta_launchpad/provider"
	"time"
)

type subscribeShare struct {
}

var SubscribeShare = new(subscribeShare)

func (s *subscribeShare) GetSubscrubeShare(userId string, subscribeId int) (ret *model.SubscribeShare, err error) {
	err = g.DB().Model("subscribe_share").Where("user_id = ? AND subscribe_id = ?", userId, subscribeId).Scan(&ret)
	if err != nil {
		return
	}

	return
}

func (s *subscribeShare) UploadSubscrubeShare(req model.SubscribeShareUpload, userId string) (err error) {
	now := time.Now()
	var ret *model.SubscribeActivity
	m := g.DB().Model("subscribe_activity")
	err = m.Where("start_time <= ? AND alias = ?", now, req.Alias).Scan(&ret)
	if err != nil {
		return
	}
	if ret == nil {
		err = fmt.Errorf("活动不存在")
		return
	}
	// 活动前一个小时
	anHourAgo := ret.ActivityEndTime.Add(-time.Hour).Time
	if now.Before(ret.ActivityStartTime.Time) {
		err = gerror.New("分享有奖活动还未开始")
		return
	}
	if now.After(anHourAgo) {
		err = gerror.New("分享有奖活动已结束")
		return
	}
	var share *model.SubscribeShare
	err = g.DB().Model("subscribe_share").Where("user_id = ? AND subscribe_id = ?", userId, ret.Id).Scan(&share)
	if err != nil {
		return
	}
	if share != nil && share.Id != 0 {
		err = fmt.Errorf("已经参与过活动")
		return
	}
	// 判断此活动是否是根据余额认购
	if ret.GeneralNumMethod == 1 {
		// 查找此用户是否有流水记录
		order, e := provider.User.GetCnyPublisherOrder(userId, req.PublisherId)
		if e != nil {
			err = e
			return
		}
		if order == nil {
			return
		}
		share = &model.SubscribeShare{
			UserId:      userId,
			SubscribeId: ret.Id,
		}
		_, err = g.DB().Model("subscribe_share").Insert(share)
		return
	}
	return
}
