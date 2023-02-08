package main

import (
	"BeeBlog/database"
	_ "BeeBlog/routers"

	_ "BeeBlog/models"

	"github.com/beego/beego/v2/adapter/orm"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/beego/beego/v2/server/web/session/mysql"
)

func main() {
	if web.BConfig.RunMode == "dev" {
		orm.Debug = true
	}
	//session设置
	web.BConfig.WebConfig.Session.SessionOn = true
	web.BConfig.WebConfig.Session.SessionProvider = "mysql"
	web.BConfig.WebConfig.Session.SessionProviderConfig = database.GetDataSource()
	web.BConfig.WebConfig.Session.SessionGCMaxLifetime = 60 * 60 * 24 * 15
	web.BConfig.WebConfig.Session.SessionCookieLifeTime = 60 * 60 * 24 * 15
	orm.RunCommand()
	//添加自定义模板函数，用户tag转换
	web.AddFuncMap("tag", database.GetTagName)
	web.Run()
}
