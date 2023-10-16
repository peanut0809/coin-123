package task

import (
	"github.com/gogf/gf/frame/g"
	"time"
)

func Banner() {
	// todo 具体逻辑
	//g.Log().Info("<=================================================>")
	StateOpen()
	StateOff()
	// 每1分钟执行一次
	time.AfterFunc(time.Minute, Banner)
}

func StateOpen() {
	type s struct {
		State int `json:"state"`
	}
	var state s
	state.State = 1
	_, err := g.DB().Model("banner").Where(" timing_state = 1 AND goods_on_time <= NOW() AND NOW() <= DATE_ADD(goods_on_time, INTERVAL +1 minute) AND state = 0 ").Update(&state)
	if err != nil {
		return
	}
}

func StateOff() {
	type s struct {
		State       int `json:"state"`
		TimingState int `json:"timing_state"`
	}
	var state s
	state.State = 0
	state.TimingState = 0
	_, err := g.DB().Model("banner").Where(" timing_state = 1 AND NOW() >= goods_off_time AND state = 1 ").Update(&state)
	if err != nil {
		return
	}
}
