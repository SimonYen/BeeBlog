package controllers

import (
	"BeeBlog/database"
	"BeeBlog/models"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

type UserController struct {
	web.Controller
}

// 主页
func (u *UserController) Home() {
	//读取flash
	web.ReadFromRequest(&u.Controller)
	//读取session，看用户是否登录过
	name := u.GetSession("user_name")
	if name != nil {
		u.Data["name"] = name.(string)
	}
	var posts []*models.Post
	qs := database.Handler.QueryTable("post")
	_, err := qs.OrderBy("-created").All(&posts)
	if err != nil {
		logs.Error(err)
	}
	u.Data["posts"] = posts
	u.Layout = "layout/base.html"
	u.TplName = "home.html"
}

// 注册用户
func (u *UserController) Register() {
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
		flash.Error(user.Email + "邮箱已被注册！")
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
		flash.Success("注册成功！")
	}
	//保存flash
	flash.Store(&u.Controller)
	//转到主页
	u.Redirect(web.URLFor("UserController.Home"), 302)
}

// 用户登录
func (u *UserController) Login() {
	flash := web.NewFlash()
	user := new(models.User)
	qs := database.Handler.QueryTable(user)
	err := u.ParseForm(user)
	if err != nil {
		logs.Error("表单解析错误：", err)
	}
	if qs.Filter("email", user.Email).Exist() {
		//比较密码是否正确
		user_in_database := new(models.User)
		qs.Filter("email", user.Email).One(user_in_database)
		if user_in_database.CheckPasswordHash(user.Password) {
			//填入到session
			u.SetSession("user_name", user_in_database.Name)
			u.SetSession("user_id", user_in_database.Id)
			flash.Success("登录成功！")
		} else {
			flash.Error("密码错误！")
		}
	} else {
		flash.Error(user.Email + "邮箱未被注册！")
	}
	flash.Store(&u.Controller)
	u.Redirect(web.URLFor("UserController.Home"), 302)
}

// 用户退出登录
func (u *UserController) Logout() {
	err := u.DestroySession()
	if err != nil {
		logs.Error(err)
	}
	u.Redirect(web.URLFor("UserController.Home"), 302)
}
