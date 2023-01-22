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
	flash := web.NewFlash()
	user := new(models.User)
	qs := database.Handler.QueryTable(user)
	err := receiver.ParseForm(user)
	if err != nil {
		logs.Error("表单解析错误：", err)
	}
	if qs.Filter("email", user.Email).Exist() {
		flash.Error(user.Email + "邮箱已被注册！🫠")
	} else {
		err = user.HashPassword()
		if err != nil {
			logs.Error("密码加密失败：", err)
		}
		//保存数据库
		_, err = database.Handler.Insert(user)
		if err != nil {
			logs.Error("数据库插入失败！", err)
		}
		flash.Success("注册成功！😉")
	}
	flash.Store(&receiver.Controller)
	receiver.Redirect("/", 302)
}
