package service

import (
	"fmt"
	"meta_launchpad/model"

	"github.com/gogf/gf/frame/g"
)

type subscribeShare struct {
}

var SubscribeShare = new(subscribeShare)

func (s *subscribeShare) GetSubscrubeShare(userId string, subscribeId int) (ret model.SubscribeShare, err error) {
	err = g.DB().Model("subscribe_share").Where("user_id = ? AND subscribe_id = ?", userId, subscribeId).Scan(&ret)
	if err != nil {
		return
	}

	return
}

func (s *subscribeShare) UploadSubscrubeShare(req model.SubscribeShareUpload, userId string) (err error) {
	var ret *model.SubscribeActivity
	m := g.DB().Model("subscribe_activity")
	err = m.Where("alias = ?", req.Alias).Scan(&ret)
	if err != nil {
		return
	}
	if ret == nil {
		err = fmt.Errorf("活动不存在")
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
	share = &model.SubscribeShare{
		UserId:      userId,
		SubscribeId: ret.Id,
	}
	_, err = g.DB().Model("subscribe_share").Insert(share)

	return
}
