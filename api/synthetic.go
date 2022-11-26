package api

import (
	"meta_launchpad/model"
	"meta_launchpad/service"
	"time"

	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/api"
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/client"
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/common/utils"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
)

type synthetic struct {
	api.CommonBase
}

var Synthetic = new(synthetic)

func (s *synthetic) Create(r *ghttp.Request) {
	var req model.SyntheticActivityReq
	err := r.Parse(&req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	req.PublisherId = s.GetPublisherId(r)
	if req.StartTime == nil || req.EndTime == nil || req.Cover == "" || req.Name == "" || req.AppId == "" || req.TemplateId == "" || req.Sum <= 0 || req.OutNum <= 0 || len(req.ConditionArr) == 0 {
		s.FailJsonExit(r, "参数错误")
		return
	}
	if req.StartTime.Unix() >= req.EndTime.Unix() {
		s.FailJsonExit(r, "结束时间必须开始时间")
		return
	}
	conditionStr := gconv.String(req.ConditionArr)
	req.SyntheticActivity.RemainNum = req.SyntheticActivity.Sum
	req.SyntheticActivity.Condition = &conditionStr
	if req.Id == 0 {
		err = service.Synthetic.Create(req.SyntheticActivity)
		if err != nil {
			s.FailJsonExit(r, err.Error())
			return
		}
		s.SusJsonExit(r)
		return
	}
	err = service.Synthetic.Update(req.SyntheticActivity)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r)
	return
}

func (s *synthetic) List(r *ghttp.Request) {
	publisherId := s.GetPublisherId(r)
	pageNum := r.GetQueryInt("pageNum", 1)
	pageSize := r.GetQueryInt("pageSize", 1)
	startTimeBegin := r.GetQueryString("startTimeBegin")
	startTimeEnd := r.GetQueryString("startTimeEnd")
	endTimeBegin := r.GetQueryString("endTimeBegin")
	endTimeEnd := r.GetQueryString("endTimeEnd")
	status := r.GetQueryString("status")
	searchVal := r.GetQueryString("searchVal")
	ret, err := service.Synthetic.List(publisherId, pageNum, pageSize, startTimeBegin, startTimeEnd, endTimeBegin, endTimeEnd, status, searchVal, "id DESC", 0, false)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
	return
}

func (s *synthetic) ClientList(r *ghttp.Request) {
	publisherId := s.GetPublisherId(r)
	pageNum := r.GetQueryInt("pageNum", 1)
	pageSize := r.GetQueryInt("pageSize", 1)
	status := r.GetQueryString("status")
	ret, err := service.Synthetic.List(publisherId, pageNum, pageSize, "", "", "", "", status, "", "start_time DESC", 1, true)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
	return
}

func (s *synthetic) Detail(r *ghttp.Request) {
	id := r.GetQueryInt("id")
	ret, err := service.Synthetic.Detail(id)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
	return
}

func (s *synthetic) Open(r *ghttp.Request) {
	id := r.GetInt("id")
	open := r.GetInt("open")
	ret, err := service.Synthetic.Open(id, open)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
	return
}

func (s *synthetic) Delete(r *ghttp.Request) {
	id := r.GetInt("id")
	ret, err := service.Synthetic.Delete(id)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
	return
}

func (s *synthetic) ClientDetail(r *ghttp.Request) {
	id := r.GetInt("id")
	ret, err := service.Synthetic.ClientDetail(id)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *synthetic) GetDoSyntheticResult(r *ghttp.Request) {
	orderNo := r.GetQueryString("orderNo")
	ret := service.Synthetic.GetResult(orderNo)
	s.SusJsonExit(r, ret)
}

func (s *synthetic) GetRecordDetail(r *ghttp.Request) {
	orderNo := r.GetQueryString("orderNo")
	aid := r.GetQueryInt("aid")
	ret, err := service.Synthetic.GetRecordList(1, 1, s.GetPublisherId(r), s.GetUserId(r), orderNo, aid)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	if len(ret.List) == 0 {
		s.FailJsonExit(r, "记录不存在")
		return
	}
	s.SusJsonExit(r, ret.List[0])
}

func (s *synthetic) GetSyntheticRecord(r *ghttp.Request) {
	pageNum := r.GetQueryInt("pageNum", 1)
	pageSize := r.GetQueryInt("pageSize", 20)
	userId := r.GetQueryString("userId")
	aid := r.GetQueryInt("aid")
	publisherId := s.GetPublisherId(r)
	ret, err := service.Synthetic.GetRecordList(pageNum, pageSize, publisherId, userId, "", aid)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *synthetic) GetRecordList(r *ghttp.Request) {
	pageNum := r.GetQueryInt("pageNum", 1)
	pageSize := r.GetQueryInt("pageSize", 20)
	userId := s.GetUserId(r)
	publisherId := s.GetPublisherId(r)
	aid := r.GetQueryInt("aid")
	ret, err := service.Synthetic.GetRecordList(pageNum, pageSize, publisherId, userId, "", aid)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, ret)
}

func (s *synthetic) DoSynthetic(r *ghttp.Request) {
	var req model.SyntheticReq
	err := r.Parse(&req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	req.UserId = s.GetUserId(r)
	// if req.UserId != "j2uuiu0pr2bcinpf5dugzvxdxhlpa2yg" {
	// 	s.FailJsonExit(r, "功能维护中!暂时无法参与!")
	// 	return
	// }
	req.OrderNo = utils.Generate()
	req.PublisherId = s.GetPublisherId(r)
	if req.UserId == "" || req.Aid == 0 {
		s.FailJsonExit(r, "参数错误!")
		return
	}
	sa, err := service.Synthetic.Detail(req.Aid)
	if err != nil {
		return
	}
	now := time.Now()
	if now.Unix() >= sa.StartTime.Unix() && now.Unix() <= sa.EndTime.Unix() {
		// ret.StatusTxt = "进行中"
	} else {
		if now.Unix() <= sa.StartTime.Unix() {
			// ret.StatusTxt = "未开始"
			s.FailJsonExit(r, "活动未开始!")
			return
		} else {
			// ret.StatusTxt = "已结束"
			s.FailJsonExit(r, "活动已结束!")
			return
		}
	}

	queueName := "synthetic.do"
	mqClient, err := client.GetQueue(client.QueueConfig{
		QueueName:  queueName,
		Exchange:   queueName,
		RoutingKey: "",
		Kind:       "fanout",
		MqUrl:      g.Cfg().GetString("rabbitmq.default.link"),
	})
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	defer mqClient.Close()
	err = mqClient.Publish(req)
	if err != nil {
		s.FailJsonExit(r, err.Error())
		return
	}
	s.SusJsonExit(r, req.OrderNo)
}
