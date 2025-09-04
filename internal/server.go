package internal

import (
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/go-http/pkg/httpserver"
	"go-wx-api/internal/api"
	"go-wx-api/internal/config"
	"go-wx-api/internal/ntfy"
)

func Bootstrap(cfg *config.Config) {
	go ntfy.GetInstance().Start(&ntfy.NtfyInfo{
		Address:  "http://v.uuxia.cn:90",
		Topic:    "uclient",
		Username: "admin",
		Password: "het002402",
	})
	ntfy.GetInstance().AddFunc(func(s string) {
		glog.Debug(s)
	})

	server := httpserver.New().
		CORSMethodMiddleware().
		AddApi(api.NewApi(cfg)).
		Done(cfg.ServerPort)
	defer server.Stop()
	server.Wait()
}
