package controllers

import (
	"BeeBlog/database"
	"BeeBlog/models"

	"github.com/beego/beego/v2/server/web"
)

type RegisterController struct {
	web.Controller
}

func (receiver *RegisterController) Post() {
	//先试着看能否收到数据
	user := new(models.User)
	err := receiver.ParseForm(user)
	if err != nil {
		println(err)
	}
	user.HashPassword()
	//保存数据库
	database.Handler.Insert(user)
	receiver.Redirect("/", 302)
}
