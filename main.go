package main

import (
	"cccn-zxl-server/router"

	"github.com/gogf/gf/frame/g"
)

func main() {
	s := router.InitRouter()
	g.Cfg().SetPath("./")

	s.Run()
}
