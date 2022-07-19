package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func SendRpc(addr, rpcName string, params interface{}) {
	var result interface{}
	d, _ := client.NewPeer2PeerDiscovery("tcp@"+addr, "")
	rpc := strings.Split(rpcName, ".")
	o := client.DefaultOption
	o.SerializeType = protocol.JSON
	client := client.NewXClient(rpc[0], client.Failtry, client.RandomSelect, d, o)
	defer client.Close()
	err := client.Call(context.Background(), rpc[1], params, &result)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("请求结果：%v \n", result)
	j, _ := json.Marshal(result)
	//j, _ := gjson.LoadContent(result)
	//jstr, _ := j.ToJsonString()
	log.Printf("请求结果：%s \n", string(j))
	ioutil.WriteFile(rpcName+".txt", j, os.ModePerm)
	return
}

func main() {
	fmt.Println(11122)
	//	time.Sleep(time.Second * 3)
	addr := "127.0.0.1:18126"
	//addr := "39.107.72.102:18121"

	//params := &map[string]interface{}{
	//	"userIds": []string{"eh2bzu01yywcj78kup54rnz2006gildw"},
	//}
	//SendRpc(addr, "UserBase.GetUserInfo", params)

	//params := &map[string]interface{}{
	//	//"userIds":  []string{"eh2bzu01yywcj78kup54rnz2006gildw"},
	//	"nickname": "小",
	//	"pageSize": 2,
	//	//"phones": []string{"13"},
	//}
	//SendRpc(addr, "UserBase.GetUserInfoByPage", params)

	//params := &map[string]interface{}{
	//	"appId":    "test",
	//	"nickname": "test",
	//	"avatar":   "test",
	//	"phone":    "13720009841",
	//	"unionId":  "unionId",
	//	"openid":   "openid",
	//}
	//SendRpc(addr, "UserBase.RegistWithWechat", params)

	//item := make(map[string]string)
	//item["xiyoudu"] = "SR"
	//params := &map[string]interface{}{
	//	"dataCondition": []map[string]string{item},
	//	"appId":         "meta_world_id",
	//	"templateId":    "41000049",
	//	"assetType":     "290",
	//	"num":           1,
	//	"userId":        "eh2bzu01yywcj78kup54rnz2006gildw",
	//}
	//SendRpc(addr, "Asset.PublishAssetWithCondition", params)

	//params := &map[string]interface{}{
	//	"appIds":   []string{"testAppId"},
	//	"tokenIds": []string{"101000000008"},
	//}
	//SendRpc(addr, "Asset.GetMateDataByAks", params)

	params := &map[string]interface{}{
		"ids": []int{89},
	}
	SendRpc(addr, "Launchpad.GetDetailByIds", params)

	//params := &map[string]interface{}{
	//	"appId":   "testAppId",
	//	"tokenId": "101000000008",
	//}
	//SendRpc(addr, "Asset.GetDetail", params)

}
