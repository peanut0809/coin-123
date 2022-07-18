package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/rpc/jsonrpc"
	"os"
	"time"
)

func SendRpc(addr, rpcName string, params *map[string]interface{}) {
	conn, e := net.DialTimeout("tcp", addr, 30*time.Second) // 30秒超时时间
	if e != nil {
		fmt.Println("服务连接失败===========" + addr)
		return
	}
	defer conn.Close()
	client := jsonrpc.NewClient(conn)
	if client != nil {
		defer client.Close()
	}
	var result interface{}
	e = client.Call(rpcName, params, &result)
	if e != nil {
		fmt.Println("服务出错===========", e)
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
	addr := "127.0.0.1:18121"
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
		"appIds":      []string{"testAppId"},
		"templateIds": []string{"101000000008"},
	}
	SendRpc(addr, "Asset.GetMateDataByTemps", params)

	//params := &map[string]interface{}{
	//	"appId":   "testAppId",
	//	"tokenId": "101000000008",
	//}
	//SendRpc(addr, "Asset.GetDetail", params)

}
