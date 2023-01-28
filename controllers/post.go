package controllers

import (
	"BeeBlog/database"
	"BeeBlog/models"
	"strconv"

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
func (p *PostController) Save() {
	//读取flash
	web.ReadFromRequest(&p.Controller)
	flash := web.NewFlash()
	//读取session，看用户是否登录过
	name := p.GetSession("user_name")
	if name != nil {
		p.Data["name"] = name.(string)
	} else {
		flash.Error("请先登录！")
		flash.Store(&p.Controller)
		p.Redirect(web.URLFor("UserController.Home"), 302)
		return
	}
	post := new(models.Post)
	err := p.ParseForm(post)
	if err != nil {
		logs.Error("表单解析失败：", err)
	}
	user := new(models.User)
	qs := database.Handler.QueryTable(user)
	qs.Filter("id", p.GetSession("user_id").(int)).One(user)
	post.Author = user
	_, err = database.Handler.Insert(post)
	if err != nil {
		logs.Error(err)
	}
	flash.Success("文章创建成功！")
	flash.Store(&p.Controller)
	p.Redirect(web.URLFor("UserController.Home"), 302)
}

// 查看文章详情
func (p *PostController) Detail() {
	//读取flash
	web.ReadFromRequest(&p.Controller)
	flash := web.NewFlash()
	//读取session，看用户是否登录过
	name := p.GetSession("user_name")
	if name != nil {
		p.Data["name"] = name
	}
	//获取文章id,从路由上
	post_id, _ := strconv.Atoi(p.Ctx.Input.Param(":id"))
	//先试着看能不能读取出来
	post := new(models.Post)
	qs_post := database.Handler.QueryTable(post)
	err := qs_post.Filter("id", post_id).One(post)
	if err != nil || post.Id == 0 {
		flash.Error("无法从数据库中找到该文章！")
		logs.Error(err)
		flash.Store(&p.Controller)
		p.Redirect(web.URLFor("UserController.Home"), 302) //之后应该转到404界面
	} else {
		//把文章和作者丢进模板就可以了
		post.Content = post.ToHTML()
		p.Data["post"] = post
		//查找作者
		author := new(models.User)
		qs_user := database.Handler.QueryTable(author)
		qs_user.Filter("id", post.Author.Id).One(author)
		p.Data["author"] = author
	}
	p.Layout = "layout/base.html"
	p.TplName = "post/view.html"
}
