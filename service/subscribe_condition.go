package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/model"
	"strings"
)

type subscribeCondition struct {
}

var SubscribeCondition = new(subscribeCondition)

//已条件分割资产
//func (s *subscribeCondition) SplitAssetByCondition(userId string, aid int) (err error) {
//	conditions, e := s.GetList(aid)
//	if e != nil {
//		err = e
//		return
//	}
//	sumAsset, e := s.GetConditionsAsset(userId, conditions)
//	if e != nil {
//		err = e
//		return
//	}
//	for _, condition := range conditions {
//		var assets []model.Asset
//		assets, err = s.GetOneConditionAsset(condition, sumAsset)
//		if err != nil {
//			return
//		}
//		count := len(assets)
//	}
//	return
//}

//获取自查条件
func (s *subscribeCondition) GetList(aid int) (ret []model.SubscribeCondition, err error) {
	err = g.DB().Model("subscribe_condition").Where("aid = ?", aid).Scan(&ret)
	if err != nil {
		return
	}
	return
}

func (s *subscribeCondition) GetConditionsAsset(userId string, in []model.SubscribeCondition) (ret []model.Asset, err error) {
	if len(in) == 0 {
		return
	}
	var rows *sql.Rows
	m := g.DB()
	begin := "COALESCE(ak.app_id,''),COALESCE(ak.token_id,''),COALESCE(am.asset_pic,''),COALESCE(am.icon,''),COALESCE(am.asset_name,''),COALESCE(am.template_id,''),COALESCE(am.type,''),COALESCE(am.data,'')"
	orSqlArr := make([]string, 0)
	for _, v := range in {
		sqlArr := make([]string, 0)
		sqlArr = append(sqlArr, fmt.Sprintf("ak.user_id = '%s'", userId))
		sqlArr = append(sqlArr, fmt.Sprintf("ak.app_id = '%s'", v.AppId))
		if v.AssetType != "" {
			sqlArr = append(sqlArr, fmt.Sprintf("am.type = '%s'", v.AssetType))
		}
		if v.TemplateId != "" {
			sqlArr = append(sqlArr, fmt.Sprintf("am.template_id = '%s'", v.TemplateId))
		}
		if v.MetaDataRule != "" {
			metaDataRule := make(map[string]string, 0)
			err = json.Unmarshal([]byte(v.MetaDataRule), &metaDataRule)
			if err != nil {
				return
			}
			for k, v := range metaDataRule {
				sqlArr = append(sqlArr, fmt.Sprintf("am.data ->> '$.%s' = '%s'", k, v))
			}
		}
		orSqlArr = append(orSqlArr, fmt.Sprintf("(%s)", strings.Join(sqlArr, " and ")))
	}
	rows, err = m.Query("select " + begin + " from meta_assets.asset_knapsack ak left join meta_assets.asset_metadata am on ak.app_id = am.app_id and ak.token_id = am.token_id where " + strings.Join(orSqlArr, " or ") + " order by ak.id desc")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var item model.Asset
		err = rows.Scan(&item.AppId, &item.TokenId, &item.AssetPic, &item.Icon, &item.AssetName, &item.TemplateId, &item.Type, &item.Data)
		if err != nil {
			return
		}
		ret = append(ret, item)
	}
	return
}

//资产条件分类
func (s *subscribeCondition) GetOneConditionAsset(in model.SubscribeCondition, assets []model.Asset) (ret []model.Asset, err error) {
	for _, v := range assets {
		vJson := make(map[string]interface{})
		err = json.Unmarshal([]byte(v.Data), &vJson)
		if err != nil {
			return
		}
		mJson := make(map[string]string, 0)
		if in.MetaDataRule != "" {
			err = json.Unmarshal([]byte(in.MetaDataRule), &mJson)
			if err != nil {
				return
			}
		}
		mdata := false
		for k, v := range mJson {
			if fmt.Sprintf("%v", vJson[k]) == fmt.Sprintf("%v", v) {
				mdata = true
			} else {
				mdata = false
				break
			}
		}
		if in.AppId == v.AppId && (in.AssetType == v.Type || in.AssetType == "") && (in.TemplateId == "" || in.TemplateId == v.TemplateId) && (in.MetaDataRule == "" || mdata) {
			v.Data = ""
			v.Type = ""
			ret = append(ret, v)
		}
	}
	return
}
