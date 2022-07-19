package rpc

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/aop"
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/rpcx"
	"context"
	"fmt"
	"meta_launchpad/model"
	"meta_launchpad/service"
	"reflect"
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
