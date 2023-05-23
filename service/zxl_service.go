package service

import (
	"fmt"
	"time"

	"github.com/gogf/gf/frame/g"
	zxl "github.com/zhixinlian/zxl-go-sdk/v2"
)

const ZXL_APP_ID = "230324002800001"
const ZXL_APP_KEY = "1ad8490c99ec40dc963be3b20ef8a3f7"

type zxlService struct {
}

var ZxlService = new(zxlService)

// ------ 创建至信链订单 -------
func (s *zxlService) MarktingZxlOrder(req *ZxlOrderCreatReq) (res *ZxlOrderCreatRes, err error) {
	err = s.MarktingZxlOrderAvalue(req)
	if err != nil {
		return
	}
	zxlImpl, err := s.MarktingZxlInit()
	if err != nil {
		return
	}
	var errMessage = ""
	// ------ 图片取证订单 ------
	if req.UrlType == "img" {
		orderNo, orderErr := zxlImpl.EvidenceObtainPic(req.WebUrls, req.Title, req.Remark, 10*time.Second)
		fmt.Println("------图片取证下单返回结果------")
		g.Log().Info(orderNo, orderErr, req)

		if orderErr != nil {
			errMessage = orderErr.Error()
		}
		res = &ZxlOrderCreatRes{
			EvidenceId:    req.EvidenceId,
			EvidenceOrder: orderNo,
			EvidenceErr:   errMessage,
		}
		return
	}
	// ------ 视频取证订单 ------
	if req.UrlType == "video" {

		if req.Duration > 5 {
			zocr := s.NewEvidenceObtainVideo(zxlImpl, req)
			return zocr, nil
		}

		videoOrder, videoOrderErr := zxlImpl.EvidenceObtainVideo(req.WebUrls, req.Title, req.Remark, 10*time.Second)
		fmt.Println("------zxl视频下单返回结果------")
		g.Log().Info(videoOrder, videoOrderErr, req)
		if videoOrderErr != nil {
			errMessage = videoOrderErr.Error()
		}
		res = &ZxlOrderCreatRes{
			EvidenceOrder: videoOrder,
			EvidenceErr:   errMessage,
		}
		return
	}
	return
}

// ------ 创建至信链订单结果查询 -------
func (s *zxlService) MarktingZxlOrderInfo(req *ZxlOrderSearchReq) (eid *zxl.EvIdData, err error) {
	if req.OrderNo == "" {
		err = fmt.Errorf("orderNo订单号异常")
		return
	}
	zxlImpl, err := s.MarktingZxlInit()
	if err != nil {
		return
	}
	eid, err = zxlImpl.GetEvidenceStatus(req.OrderNo, 10*time.Second)
	fmt.Println("------zxl下单返回结果------")
	g.Log().Info(eid, err, req)
	return
}

// ------ 校验zxl下单请求参数 ------
func (s *zxlService) MarktingZxlOrderAvalue(req *ZxlOrderCreatReq) (err error) {
	if req.Remark == "" {
		err = fmt.Errorf("Remark异常")
		return
	}
	if req.Title == "" {
		err = fmt.Errorf("Title异常")
		return
	}
	if req.UrlType == "" {
		err = fmt.Errorf("UrlType异常")
		return
	}

	if req.WebUrls == "" {
		err = fmt.Errorf("WebUrls异常")
		return
	}
	return
}

// ------ 初始化zxl指针 ------
func (s *zxlService) MarktingZxlInit() (zxlImpl *zxl.ZxlImpl, err error) {
	zxlImpl, err = zxl.NewZxlImpl(ZXL_APP_ID, ZXL_APP_KEY)
	fmt.Println("------ZXLInit初始化MarktingZxlInit------")
	g.Log().Info(zxlImpl, err)
	return
}

// ------ 长视频录屏 ------
func (s *zxlService) NewEvidenceObtainVideo(zxlImpl *zxl.ZxlImpl, req *ZxlOrderCreatReq) (res *ZxlOrderCreatRes) {
	fmt.Println("------长视频录屏下单-------")
	if req.Duration > 60 {
		req.Duration = 60
	}
	orderReq := zxl.ObtainVideoOption{
		WebUrls:        req.WebUrls,
		Title:          req.Title,
		Remark:         req.Remark,
		RepresentAppId: "",
		Duration:       req.Duration * 60,
	}
	fmt.Println("------zxl长视频下单请求------")
	g.Log().Info(orderReq, req)
	videoOrder, videoOrderErr := zxlImpl.NewEvidenceObtainVideo(&orderReq, 10*time.Second)
	fmt.Println("------zxl长视频下单返回结果------")
	g.Log().Info(videoOrder, videoOrderErr, req)
	errMessage := ""
	if videoOrderErr != nil {
		errMessage = videoOrderErr.Error()
	}
	res = &ZxlOrderCreatRes{
		EvidenceOrder: videoOrder,
		EvidenceErr:   errMessage,
	}
	return
}

// ------ 创建至信链订单请求参数 -------
type ZxlOrderCreatReq struct {
	Title      string `json:"title"`      // 取证标题
	Remark     string `json:"remark"`     // 取证b备注
	WebUrls    string `json:"webUrls"`    // 取证图片或者视频链接
	UrlType    string `json:"urlType"`    // url类型
	Timeout    int    `json:"timeout"`    // 超时时间
	Duration   int    `json:"duration"`   // 录屏时间
	EvidenceId int    `json:"evidenceId"` // 后台创建取证请求id
}

// ------ 创建至信链订单请求参数 -------
type ZxlOrderCreatRes struct {
	EvidenceId    int    `json:"evidenceId"`
	EvidenceOrder string `json:"evidenceOrder"`
	EvidenceErr   string `json:"evidenceErr"`
}

// ------ 创建至信链订单请求参数 -------
type ZxlOrderSearchReq struct {
	EvidenceId int    `json:"evidenceId"` // node服务那边的id
	OrderNo    string `json:"orderNo"`    // 订单号
}
