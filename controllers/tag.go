package controllers

import (
	"BeeBlog/database"
	"BeeBlog/models"

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
