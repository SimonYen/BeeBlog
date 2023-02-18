package controllers

import (
	"BeeBlog/database"
	"BeeBlog/models"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

type TagController struct {
	web.Controller
}

func (t *TagController) ChangeTag() {
	//读取flash
	web.ReadFromRequest(&t.Controller)
	flash := web.NewFlash()
	//读取session，看用户是否登录过
	user_id := t.GetSession("user_id")
	if user_id == nil {
		flash.Error("请先登录！")
		flash.Store(&t.Controller)
		t.Redirect(web.URLFor("UserController.Home"), 302)
		return
	}
	post_id, _ := t.GetInt("id")
	//查询作者
	qs_p := database.Handler.QueryTable("post")
	post := new(models.Post)
	qs_p.Filter("Id", post_id).One(post)
	//查看本地用户是否和作者一致
	if post.Author.Id != user_id.(int) {
		flash.Error("你无权修改别人的文章类别！")
		flash.Store(&t.Controller)
		t.Redirect(web.URLFor("UserController.Profile"), 302)
		return
	}
	//获取修改后的tag
	tag_id, _ := t.GetInt("tag")
	post.Class.Id = tag_id
	//保存到数据库
	database.Handler.Update(post, "Class")
	flash.Success("成功修改类别。")
	flash.Store(&t.Controller)
	t.Redirect(web.URLFor("UserController.Profile"), 302)
}

func (t *TagController) TaggedHome() {
	//读取flash
	web.ReadFromRequest(&t.Controller)
	//读取session，看用户是否登录过
	name := t.GetSession("user_name")
	if name != nil {
		t.Data["name"] = name.(string)
	}
	//获取tag id
	tag_id, err := t.GetInt(":id")
	if err != nil {
		flash := web.NewFlash()
		flash.Error("地址栏tag_id解析错误，不是整数！")
		flash.Store(&t.Controller)
		t.Redirect(t.URLFor("UserController.Home"), 302)
		return
	}
	var posts []*models.Post
	var tags []*models.Tag
	qs_p := database.Handler.QueryTable("post")
	qs_t := database.Handler.QueryTable("tag")
	qs_u := database.Handler.QueryTable("user")
	user_id_ := t.GetSession("user_id")
	user := new(models.User)
	class := new(models.Tag)
	qs_t.Filter("id", tag_id).One(class)
	if class.Id == 0 {
		flash := web.NewFlash()
		flash.Error("没有这个分类！")
		flash.Store(&t.Controller)
		t.Redirect(t.URLFor("UserController.Home"), 302)
		return
	}
	t.Data["class"] = class
	if user_id_ != nil {
		user_id := user_id_.(int)
		qs_u.Filter("id", user_id).One(user)
		t.Data["user"] = user
	}
	_, err = qs_p.Filter("Class__Id", tag_id).OrderBy("-created").All(&posts)
	if err != nil {
		logs.Error(err)
	}
	_, err = qs_t.All(&tags)
	if err != nil {
		logs.Error(err)
	}
	t.Data["posts"] = posts
	t.Data["tags"] = tags
	t.Layout = "layout/base.html"
	t.TplName = "home.html"
}
