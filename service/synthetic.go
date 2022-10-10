package service

import (
	"github.com/gogf/gf/frame/g"
	"meta_launchpad/model"
)

type synthetic struct {
}

var Synthetic = new(synthetic)

func (s *synthetic) Create(in model.SyntheticActivity) (err error) {
	m := g.DB().Model("synthetic_activity")
	_, err = m.Insert(&in)
	if err != nil {
		return
	}
	return
}

func (s *synthetic) Update(in model.SyntheticActivity) (err error) {
	//有人合成就不能修改

	m := g.DB().Model("synthetic_activity")
	_, err = m.Data(in).Where("id", in.Id).Update()
	if err != nil {
		return
	}
	return
}
