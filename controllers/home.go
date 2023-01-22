package controllers

import (
	"github.com/beego/beego/v2/server/web"
)

//主页控制器，负责主页博客渲染
type HomeController struct {
	web.Controller
}

func (receiver *HomeController) Get() {
	//添加flash消息处理
	//尝试获取session
	smsg := receiver.GetSession("smsg")
	if smsg != nil {
		receiver.Data["success_msg"] = smsg.(string)
		receiver.DelSession("smsg")
	}
	receiver.Layout = "layout/base.html"
	receiver.TplName = "home.html"
}
