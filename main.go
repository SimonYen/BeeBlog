package main

import (
	_ "BeeBlog/routers"

	_ "BeeBlog/models"

	"github.com/beego/beego/v2/adapter/orm"
	"github.com/beego/beego/v2/server/web"
)

func main() {
	if web.BConfig.RunMode == "dev" {
		orm.Debug = true
	}
	orm.RunCommand()
	web.Run()
}
