package rpc

import (
	"context"
	"fmt"
	"meta_launchpad/model"
	"meta_launchpad/service"
	"reflect"
	"time"

	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/aop"
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/rpcx"
	"github.com/gogf/gf/frame/g"
)

type Launchpad struct {
}

var LaunchpadRpc = new(Launchpad)

func init() {
	aop.RegisterAOP(reflect.TypeOf(&LaunchpadRpc))
	rpcx.RegisterRpc(LaunchpadRpc)
}

type GetDetailByIdsReq struct {
	Ids []int `json:"ids"`
}

func (t *Launchpad) GetDetailByIds(ctx context.Context, req *GetDetailByIdsReq, result *interface{}) (err error) {
	if len(req.Ids) > 100 {
		err = fmt.Errorf("参数错误")
		return
	}
	var as []model.SubscribeActivity
	as, err = service.SubscribeActivity.GetListSimple(req.Ids)
	if err != nil {
		return
	}
	*result = as
	return
}

// 根据template_id 获取活动状态
type GetEquityByTemplateIdsReq struct {
	TemplateIds []string `json:"templateIds"`
}

func (t *Launchpad) GetEquityByTemplateIds(ctx context.Context, req *GetEquityByTemplateIdsReq, result *interface{}) (err error) {
	var equityActivity []model.EquityActivity
	timeCrv := time.Now().Unix()
	nowTime := time.Unix(timeCrv, 0).Format("2006-01-02 15:04:05")
	err = g.DB().Model("equity_activity").Where("template_id IN (?)", req.TemplateIds).Where("status", 1).Where("activity_end_time > ", nowTime).Scan(&equityActivity)
	if err != nil {
		return
	}
	equityInfoMap := make(map[string]model.EquityActivity)
	for _, v := range equityActivity {
		equityInfoMap[v.TemplateId] = v
	}
	*result = equityInfoMap
	return
}

func (t *Launchpad) OffShelvesEquityActivity(ctx context.Context, req interface{}, result *interface{}) (err error) {
	err = service.AdminEquity.OffShelvesEquityActivity()
	return
}
