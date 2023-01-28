package controllers

import (
	"BeeBlog/database"
	"BeeBlog/models"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

type PostController struct {
	web.Controller
}

// 新增文章界面显示
func (p *PostController) Create() {
	//读取flash
	web.ReadFromRequest(&p.Controller)
	//读取session，看用户是否登录过
	name := p.GetSession("user_name")
	if name != nil {
		p.Data["name"] = name.(string)
	} else {
		flash := web.NewFlash()
		flash.Error("请先登录！")
		flash.Store(&p.Controller)
		p.Redirect(web.URLFor("HomeController.Get"), 302)
		return
	}
	p.Layout = "layout/base.html"
	p.TplName = "post/editor.html"
}

// 保存文章
func (receiver *PostController) Save() {
	//读取flash
	web.ReadFromRequest(&receiver.Controller)
	flash := web.NewFlash()
	//读取session，看用户是否登录过
	name := receiver.GetSession("user_name")
	if name != nil {
		receiver.Data["name"] = name.(string)
	} else {
		flash.Error("请先登录！")
		flash.Store(&receiver.Controller)
		receiver.Redirect(web.URLFor("UserController.Home"), 302)
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
	flash.Success("文章创建成功！")
	flash.Store(&receiver.Controller)
	receiver.Redirect(web.URLFor("UserController.Home"), 302)
}
