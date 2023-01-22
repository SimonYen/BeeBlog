package controllers

import (
	"BeeBlog/database"
	"BeeBlog/models"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

type RegisterController struct {
	web.Controller
}

func (receiver *RegisterController) Post() {
	//先试着看能否收到数据
	user := new(models.User)
	err := receiver.ParseForm(user) //todo 添加数据校验
	if err != nil {
		logs.Error("表单解析错误：", err)
	}
	err = user.HashPassword()
	if err != nil {
		logs.Error("密码加密失败：", err)
	}
	//保存数据库
	database.Handler.Insert(user)
	err = receiver.SetSession("smsg", "注册成功！")
	if err != nil {
		logs.Error("session保存失败：", err)
	}
	receiver.Redirect("/", 302)
}
