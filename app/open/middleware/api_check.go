package middleware

import (
	"brq5j1d.gfanx.pro/meta_cloud/meta_common/library"

	"brq5j1d.gfanx.pro/meta_cloud/meta_service/app/third/codes"
	"github.com/gogf/gf/errors/gerror"

	"github.com/gogf/gf/encoding/gjson"

	"github.com/gogf/gf/net/ghttp"
)

func ApiCheck(r *ghttp.Request) {
	_, err := GetAppId(r.GetBody())
	if err != nil {
		library.FailJsonCodeExit(r, err)
		return
	}

	r.Middleware.Next()
}

func GetAppId(data []byte) (appId string, err error) {
	reqMap := make(map[string]string)
	// err = json.Unmarshal(data, &reqMap)
	// if err != nil {
	// 	err = gerror.NewCode(codes.ParamError)
	// 	return
	// }
	j := gjson.New(data)
	j.Struct(&reqMap)

	appId = reqMap["appId"]
	if appId == "" {
		err = gerror.NewCode(codes.ParamError)
	}

	return
}
