package main

import (
	"peanut-coin123/router"

	"github.com/gogf/gf/frame/g"
)

func main() {
	s := router.InitRouter()
	g.Cfg().SetPath("./")

	s.Run()
}
