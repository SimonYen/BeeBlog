package controllers

import "github.com/beego/beego/v2/server/web"

type PostAddController struct {
	web.Controller
}

// 新增文章界面显示
func (receiver *PostAddController) Get() {
	//读取flash
	web.ReadFromRequest(&receiver.Controller)
	//读取session，看用户是否登录过
	name := receiver.GetSession("user_name")
	if name != nil {
		receiver.Data["name"] = name.(string)
	}
	receiver.Layout = "layout/base.html"
	receiver.TplName = "post/editor.html"
}
