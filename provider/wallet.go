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

// WalletAuthentication 钱包认证状态
type WalletAuthentication struct {
	IsActivation int    `json:"is_activation"` //激活状态,0 - 未认证 1 - 已认证 2：处理中 3：被驳回
	VerifyStatus string `json:"verify_status"` //进度状态
	Account      int    `json:"account"`       //钱包余额
	Count        int    `json:"count"`         //银行卡个数
	Msg          string `json:"msg"`           //驳回信息
}

// WalletAuthenticationState 钱包认证状态
func (s *wallet) WalletAuthenticationState(userId string) (wallerRet *WalletAuthentication, err error) {
	walletParams := map[string]interface{}{
		"out_user_id": userId,
	}
	err = utils.SendJsonRpcScan(context.Background(), "mate-wallet", "Account.Activation", walletParams, &wallerRet)
	if err != nil {
		g.Log().Errorf("rpc Account.Activation err:%v", err)
		return
	}
	return
}
