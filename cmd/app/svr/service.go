package svr

import (
	"fmt"
	"github.com/kardianos/service"
	"github.com/xxl6097/glog/glog"
	"github.com/xxl6097/go-service/pkg/gs/igs"
	"github.com/xxl6097/go-service/pkg/ukey"
	"github.com/xxl6097/go-service/pkg/utils"
	"github.com/xxl6097/go-wx-api/internal"
	"github.com/xxl6097/go-wx-api/internal/config"
	"github.com/xxl6097/go-wx-api/pkg"
	"os"
)

type Service struct {
	timestamp string
	gs        igs.Service
}

func (t *Service) OnStop() {
	glog.Info("service stop")
}

func (t *Service) OnShutdown() {
	glog.Info("OnShutdown ...")
}

func (this *Service) OnFinish() {
}

func load() (*config.Config, error) {
	defer glog.Flush()
	byteArray, err := ukey.Load()
	if err != nil {
		return nil, err
	}
	var cfg config.Config
	err = ukey.GobToStruct(byteArray, &cfg)
	//err = json.Unmarshal(byteArray, &cfg)
	if err != nil {
		glog.Println("ClientConfig解析错误", err)
		return nil, err
	}
	pkg.Version()
	return &cfg, nil
}

func (this *Service) OnConfig() *service.Config {
	cfg := service.Config{
		Name: pkg.AppName,
		//UserName:    "root",
		DisplayName: fmt.Sprintf("%s_%s", pkg.AppName, pkg.AppVersion),
		Description: "A Golang AAATest Service..",
	}
	return &cfg
}

func (this *Service) OnVersion() string {
	pkg.Version()
	cfg, err := load()
	if err == nil {
		glog.Debugf("cfg:%+v", cfg)
	}
	return pkg.AppVersion
}

func (this *Service) OnRun(service igs.Service) error {
	this.gs = service
	cfg, err := load()
	if err != nil {
		return err
	}
	glog.Debug("程序运行", os.Args)
	internal.Bootstrap(cfg)
	return nil
}

func (this *Service) GetAny(s2 string) []byte {
	return this.menu()
}

func (this *Service) menu() []byte {
	port := utils.InputIntDefault(fmt.Sprintf("输入服务端口(%d)：", 80), 80)
	username := utils.InputStringEmpty(fmt.Sprintf("输入管理用户名(%s)：", "admin"), "admin")
	password := utils.InputStringEmpty(fmt.Sprintf("输入管理密码(%s)：", "admin"), "admin")
	cfg := &config.Config{ServerPort: port, Username: username, Password: password, AppID: "wxbe2c2961b236427f", AppSecret: "667fc391b1ca8f4c58d1b5f224356ad5"}
	bb, e := ukey.StructToGob(cfg)
	if e != nil {
		return nil
	}
	return bb
}
