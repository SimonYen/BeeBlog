package controllers

import (
	"BeeBlog/database"
	"BeeBlog/models"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

/*
看了官方最新教程，开始重构
*/

type UserController struct {
	web.Controller
}

func (u UserController) CreateUser() {
	//初始化flash
	flash := web.NewFlash()
	user := new(models.User)
	qs := database.Handler.QueryTable(user)
	//解析表单
	err := u.ParseForm(user)
	if err != nil {
		logs.Error(err)
	}
	//检查邮箱是否存在
	if qs.Filter("email", user.Email).Exist() {
		flash.Error(user.Email + "邮箱已被注册！🫠")
	} else {
		//哈希加密
		err = user.HashPassword()
		if err != nil {
			logs.Error(err)
		}
		//保存到数据库
		_, err = database.Handler.Insert(user)
		if err != nil {
			logs.Error(err)
		}
		flash.Success("注册成功！😉")
	}
	//保存flash
	flash.Store(&u.Controller)
	//转到主页
	u.Redirect(web.URLFor(""), 302)
}
