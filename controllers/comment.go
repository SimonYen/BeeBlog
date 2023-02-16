package controllers

import (
	"BeeBlog/database"
	"BeeBlog/models"
	"fmt"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

type CommentController struct {
	web.Controller
}

// 添加评论
func (c *CommentController) Add() {
	//读取flash
	web.ReadFromRequest(&c.Controller)
	flash := web.NewFlash()
	//读取session，看用户是否登录过
	name := c.GetSession("user_name")
	if name == nil {
		flash.Error("请先登录！")
		flash.Store(&c.Controller)
		c.Redirect(c.URLFor("UserController.Home"), 302)
		return
	}
	c.Data["name"] = name.(string)
	//获取用户id
	user_id := c.GetSession("user_id")
	//从url中获取文章id
	post_id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	comment := new(models.Comment)
	c.ParseForm(comment)
	user := new(models.User)
	user.Id = user_id.(int)
	post := new(models.Post)
	post.Id = post_id
	comment.Belong = post
	comment.Author = user
	//插入到数据库中
	database.Handler.Insert(comment)
	flash.Success("评论成功。")
	flash.Store(&c.Controller)
	url := fmt.Sprintf("%s/%d", "/post", post_id)
	//urlfor还是不会用，beego的文档也太他妈垃圾了吧
	c.Redirect(url, 302)
}

func (c *CommentController) Delete() {
	//读取flash
	web.ReadFromRequest(&c.Controller)
	flash := web.NewFlash()
	//读取session，看用户是否登录过
	name := c.GetSession("user_name")
	id := c.GetSession("user_id")
	if name != nil {
		c.Data["name"] = name
	}
	//获取评论id,从路由上
	comment_id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	//先试着看能不能读取出来
	comment := new(models.Comment)
	qs_comment := database.Handler.QueryTable("comment")
	err := qs_comment.Filter("id", comment_id).One(comment)
	if err != nil || comment.Id == 0 {
		flash.Error("无法从数据库中找到该评论！")
		logs.Error(err)
		flash.Store(&c.Controller)
		c.Redirect(web.URLFor("UserController.Profile"), 302) //之后应该转到404界面
		return
	} else if comment.Author.Id != id.(int) {
		flash.Error("无权删除别人的评论！")
		flash.Store(&c.Controller)
		c.Redirect(web.URLFor("UserController.Profile"), 302)
		return

	} else {
		//删除文章
		database.Handler.Delete(comment)
		flash.Success("评论删除成功。")
		flash.Store(&c.Controller)
		c.Redirect(web.URLFor("UserController.Profile"), 302)
	}
}
