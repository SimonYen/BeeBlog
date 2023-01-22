package main

import (
	_ "BeeBlog/routers"

	_ "BeeBlog/models"

	"github.com/beego/beego/v2/adapter/orm"
	"github.com/beego/beego/v2/server/web"
	// _ "github.com/beego/beego/v2/server/web/session/redis"
)

func main() {
	if web.BConfig.RunMode == "dev" {
		orm.Debug = true
	}
	//session设置
	web.BConfig.WebConfig.Session.SessionOn = true
	// web.BConfig.WebConfig.Session.SessionProvider = "redis"
	// web.BConfig.WebConfig.Session.SessionProviderConfig = "127.0.0.1:6379"
	orm.RunCommand()
	web.Run()
}
