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
	tag := new(models.Tag)
	post.Class = tag
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
		user_id := p.GetSession("user_id").(int)
		user := new(models.User)
		qs_user.Filter("id", user_id).One(user)
		p.Data["user"] = user
		p.Data["author"] = author
		//查找评论
		var comments []*models.Comment
		qs_comment := database.Handler.QueryTable("comment")
		qs_comment.Filter("Belong__Id", post_id).All(&comments)
		//接着再将用户名和头像填入进去
		for _, comment := range comments {
			tmp := new(models.User)
			qs_user.Filter("id", comment.Author.Id).One(tmp)
			comment.Author = tmp
		}
		p.Data["comments"] = comments
	}
	p.Layout = "layout/base.html"
	p.TplName = "post/view.html"
}

// 删除文章，注意需要事先比较下作者是否是本机登录用户
func (p *PostController) Delete() {
	//读取flash
	web.ReadFromRequest(&p.Controller)
	flash := web.NewFlash()
	//读取session，看用户是否登录过
	name := p.GetSession("user_name")
	id := p.GetSession("user_id")
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
		return
	} else if post.Author.Id != id.(int) {
		flash.Error("无权删除别人的文章！")
		flash.Store(&p.Controller)
		p.Redirect(web.URLFor("UserController.Profile"), 302)
		return

	} else {
		//删除文章
		database.Handler.Delete(post)
		flash.Success("文章删除成功。")
		flash.Store(&p.Controller)
		p.Redirect(web.URLFor("UserController.Profile"), 302)
	}
}
