package common

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type commonMeans struct {
}

var CommonMeans = new(commonMeans)

// ------ 返回成功 ------
func (c *commonMeans) ResponseSuccess(r *ghttp.Request, responseData interface{}) {
	response := &g.Map{
		"code": 0,
		"msg":  "success",
		"data": responseData,
	}
	r.SetParam("apiReturnRes", response)
	r.Response.WriteJson(response)
}

// ------ 返回失败 ------
func (c *commonMeans) ResponseFail(r *ghttp.Request,message string) {
	if message == ""{
		message = "fail"
	}
	response := &g.Map{
		"code": -1,
		"msg":  message,
		"data": g.Map{},
	}
	r.SetParam("apiReturnRes", response)
	r.Response.WriteJson(response)
}
