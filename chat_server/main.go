package main

import (
	"time"

	_ "github.com/tiancai110a/chat_demo/errno"
)

var mgr *UserMgr

func main() {
	initRedis("localhost:6379", 16, 1024, time.Second*300)
	mgr = NewUserMgr(pool)
	runServer("0.0.0.0:10000")

}
