package provider

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	"context"
	"github.com/gogf/gf/frame/g"
	"strings"
)

type wallet struct {
}

var Wallet = new(wallet)

type WalletBalance struct {
	SettledAmount   int `json:"settled_amount"`   //已清算余额
	PendingAmount   int `json:"pending_amount"`   //在途余额
	ExpensingAmount int `json:"expensing_amount"` //不可用余额
}

func (s *user) GetWalletBalance(userId, acctScene string) (res *WalletBalance, err error) {
	walletBalanceParams := map[string]interface{}{
		"out_user_id": userId,
		"acct_scene":  acctScene,
	}
	err = utils.SendJsonRpcScan(context.Background(), "mate-wallet", "Scene.GetBalanceAcctsService", walletBalanceParams, &res)
	if err != nil && !strings.Contains(err.Error(), "[error]") {
		g.Log().Errorf("rpc Scene.GetBalanceAcctsServiceerr:%v", err)
		return
	}
	return
}
