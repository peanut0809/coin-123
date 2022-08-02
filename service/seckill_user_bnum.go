package service

import (
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/model"
	"strings"
)

var SeckillUserBnum = new(seckillUserBnum)

type seckillUserBnum struct {
}

func (s *seckillUserBnum) CreateAndDecr(tx *gdb.TX, in model.SeckillUserBnum, buyNum int) (err error) {
	_, err = tx.Model("seckill_user_bnum").Insert(&in)
	if err != nil {
		if strings.Contains(err.Error(), "udx_aid_uid") {
			err = nil
			r, e := tx.Exec("UPDATE seckill_user_bnum SET can_buy = can_buy - ? WHERE user_id = ? AND aid = ?", buyNum, in.UserId, in.Aid)
			if e != nil {
				err = fmt.Errorf("超过了最大购买数量")
				return
			}
			affectedNum, _ := r.RowsAffected()
			if affectedNum != 1 {
				err = fmt.Errorf("内部错误")
				return
			}
			return
		}
		return
	}
	r, e := tx.Exec("UPDATE seckill_user_bnum SET can_buy = can_buy - ? WHERE user_id = ? AND aid = ?", buyNum, in.UserId, in.Aid)
	if e != nil {
		err = fmt.Errorf("超过了单人可购买最大数量")
		return
	}
	affectedNum, _ := r.RowsAffected()
	if affectedNum != 1 {
		err = fmt.Errorf("内部错误")
		return
	}
	return
}

func (s *seckillUserBnum) GetDetail(userId string, aid int) (ret *model.SeckillUserBnum, err error) {
	err = g.DB().Model("seckill_user_bnum").Where("user_id = ? AND aid = ?", userId, aid).Scan(&ret)
	return
}

func (s *seckillUserBnum) UpdateRemain(tx *gdb.TX, userId string, aid int, num int) (err error) {
	_, err = tx.Exec("UPDATE seckill_user_bnum SET can_buy = can_buy + ? WHERE user_id = ? AND aid = ?", num, userId, aid)
	return
}
