package controllers

import (
	"BeeBlog/database"
	"BeeBlog/models"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

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
	} else {
		flash := web.NewFlash()
		flash.Error("请先登录！")
		flash.Store(&receiver.Controller)
		receiver.Redirect(web.URLFor("HomeController.Get"), 302)
		return
	}
	receiver.Layout = "layout/base.html"
	receiver.TplName = "post/editor.html"
}

// 保存文章
func (receiver *PostAddController) Post() {
	//读取flash
	web.ReadFromRequest(&receiver.Controller)
	//读取session，看用户是否登录过
	name := receiver.GetSession("user_name")
	if name != nil {
		receiver.Data["name"] = name.(string)
	} else {
		flash := web.NewFlash()
		flash.Error("请先登录！")
		flash.Store(&receiver.Controller)
		receiver.Redirect(web.URLFor("HomeController.Get"), 302)
		return
	}
	post := new(models.Post)
	err := receiver.ParseForm(post)
	if err != nil {
		logs.Error("表单解析失败：", err)
	}
	user := new(models.User)
	qs := database.Handler.QueryTable(user)
	qs.Filter("id", receiver.GetSession("user_id").(int)).One(user)
	post.Author = user
	_, err = database.Handler.Insert(post)
	if err != nil {
		logs.Error(err)
	}
	receiver.Layout = "layout/base.html"
	receiver.TplName = "home.html"
}
