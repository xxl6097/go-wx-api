package main

import (
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/go-service/pkg/gs"
	_ "github.com/xxl6097/go-wx-api/assets/buffer"
	"github.com/xxl6097/go-wx-api/cmd/app/svr"
	"github.com/xxl6097/go-wx-api/internal"
	"github.com/xxl6097/go-wx-api/internal/config"
	"github.com/xxl6097/go-wx-api/internal/u"
)

func init() {
	//go func() { log.Println(http.ListenAndServe("localhost:6060", nil)) }()
	if u.IsMacOs() {
		internal.Bootstrap(&config.Config{
			Username:   "admin",
			Password:   "admin",
			ServerPort: 9091,
			AppID:      "wxbe2c2961b236427f",
			AppSecret:  "667fc391b1ca8f4c58d1b5f224356ad5",
		})
	}
}

func main() {
	s := svr.Service{}
	err := gs.Run(&s)
	glog.Debug("程序结束", err)
}
