package cache

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
)

const DISTRIBUTED_LOCK = "meta_launchpad:lock:%s"
const SUB_PAY_TIMEOUT = "meta_launchpad:sub:timeout:%s_%d"

//分布式锁
func DistributedLock(taskName string) bool {
	re, e := g.Redis().Do("SET", fmt.Sprintf(DISTRIBUTED_LOCK, taskName), 1, "ex", 3600, "nx")
	if fmt.Sprintf("%v", re) == "OK" && e == nil {
		return true
	}
	return false
}

func DistributedUnLock(taskName string) {
	g.Redis().Do("DEL", fmt.Sprintf(DISTRIBUTED_LOCK, taskName))
}
