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

type LoginController struct {
	web.Controller
}

type LogoutController struct {
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
	receiver.Redirect(web.URLFor("HomeController.Get"), 302)
}

func (receiver *LoginController) Post() {
	flash := web.NewFlash()
	user := new(models.User)
	qs := database.Handler.QueryTable(user)
	err := receiver.ParseForm(user)
	if err != nil {
		logs.Error("表单解析错误：", err)
	}
	if qs.Filter("email", user.Email).Exist() {
		//比较密码是否正确
		u := new(models.User)
		qs.Filter("email", user.Email).One(u)
		if u.CheckPasswordHash(user.Password) {
			//填入到session
			receiver.SetSession("user_name", u.Name)
			receiver.SetSession("user_id", u.Id)
			flash.Success("登录成功！😉")
		} else {
			flash.Error("密码错误！🙃")
		}
	} else {
		flash.Error(user.Email + "邮箱未被注册！🫠")
	}
	flash.Store(&receiver.Controller)
	receiver.Redirect(web.URLFor("HomeController.Get"), 302)
}

func (receiver *LogoutController) Get() {
	err := receiver.DestroySession()
	if err != nil {
		logs.Error(err)
	}
	receiver.Redirect(web.URLFor("HomeController.Get"), 302)
}
