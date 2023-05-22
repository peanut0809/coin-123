package api

import (
	"cccn-zxl-server/common"
	"cccn-zxl-server/service"
	"fmt"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	zxl "github.com/zhixinlian/zxl-go-sdk/v2"
)

type zxlApi struct {
}

var ZxlApi = new(zxlApi)

func (s *zxlApi) Demo(r *ghttp.Request) {
	fmt.Println("-----------------")
	zxl, err := zxl.NewZxlImpl("230324002800001", "1ad8490c99ec40dc963be3b20ef8a3f7")
	fmt.Println(zxl)
	// var webUrls = "https://detail.tmall.com/item.htm?ali_refid=a3_430673_1006:1107508892:N:FJo75ygmnRDhLGcUWE8Fjg==:01f3802a75d4d7d825ed9ed7473e7990&ali_trackid=1_01f3802a75d4d7d825ed9ed7473e7990&id=610896655266&skuId=4463712938759&spm=a2e0b.20350158.31919782.1"
	// var title = "尚萌 凹凸世界 正版动漫周边嘉德罗斯雷狮公仔小挂件趴趴玩偶娃娃"
	// var remark = "尚萌 凹凸世界 正版动漫周边嘉德罗斯雷狮公仔小挂件趴趴玩偶娃娃"

	fmt.Println("-----------图片1684315355296985205099436查询取证结果-------")
	eid, err2 := zxl.GetEvidenceStatus("1684315355296985205099436", 10*time.Second)
	orderRex := g.Map{
		"eid":  eid,
		"eerr": err2,
	}
	g.Log("-------zxl图片请求和结果------").Info(orderRex)

	fmt.Println("-----------视频1684315781121697056569066查询取证结果-------")
	vide, err3 := zxl.GetEvidenceStatus("1684315781121697056569066", 10*time.Second)
	orderRex2 := g.Map{
		"eid":  vide,
		"eerr": err3,
	}
	g.Log("-------zxl视频请求和结果------").Info(orderRex2)

	return

	var webUrls = "https://www.bilibili.com/bangumi/play/ep114687?spm_id_from=autoNext"
	var title = "凹凸世界第二季 淘汰赛"
	var remark = "凹凸世界第二季 淘汰赛"

	reqMap := g.Map{
		"webUrls": webUrls,
		"title":   title,
		"remark":  remark,
		"appId":   "230324002800001",
		"appKey":  "1ad8490c99ec40dc963be3b20ef8a3f7",
	}
	fmt.Println("------zxl请求参数------")
	g.Log().Info(reqMap)

	//orderNo, err := zxl.EvidenceObtainPic(webUrls, title, remark, 10*time.Second)

	orderNo, err := zxl.EvidenceObtainVideo(webUrls, title, remark, 10*time.Second)
	fmt.Println("------zxl下单返回结果------")
	g.Log().Info(orderNo)
	g.Log().Info(err)

}

// ------ 至信链创建订单------
func (s *zxlApi) MarktingZxlOrder(r *ghttp.Request) {
	req := &service.ZxlOrderCreatReq{}
	err := r.Parse(&req)
	if err != nil {
		return
	}
	zocr, err2 := service.ZxlService.MarktingZxlOrder(req)
	if err2 != nil {
		common.CommonMeans.ResponseFail(r,err2.Error())
		return
	}
	common.CommonMeans.ResponseSuccess(r, zocr)
}

// ------ 根据订单号获取订单结果主动轮询请求------
func (s *zxlApi) MarktingZxlOrderInfo(r *ghttp.Request) {
	req := &service.ZxlOrderSearchReq{}
	err := r.Parse(&req)
	if err != nil {
		return
	}
	eid, err2 := service.ZxlService.MarktingZxlOrderInfo(req)
	if err2 != nil {
		common.CommonMeans.ResponseFail(r,err2.Error())
		return
	}
	common.CommonMeans.ResponseSuccess(r, eid)
}
