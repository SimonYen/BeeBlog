package controllers

import (
	"github.com/beego/beego/v2/server/web"
)

//主页控制器，负责主页博客渲染
type HomeController struct {
	web.Controller
}

func (receiver *HomeController) Get() {
	//读取flash
	web.ReadFromRequest(&receiver.Controller)
	receiver.Layout = "layout/base.html"
	receiver.TplName = "home.html"
}
