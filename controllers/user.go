package controllers

import (
	"BeeBlog/database"
	"BeeBlog/models"
	"path"
	"strconv"

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
	var tags []*models.Tag
	qs_p := database.Handler.QueryTable("post")
	qs_t := database.Handler.QueryTable("tag")
	qs_u := database.Handler.QueryTable("user")
	user_id_ := u.GetSession("user_id")
	user := new(models.User)
	if user_id_ != nil {
		user_id := user_id_.(int)
		qs_u.Filter("id", user_id).One(user)
		u.Data["user"] = user
	}
	_, err := qs_p.OrderBy("-created").All(&posts)
	if err != nil {
		logs.Error(err)
	}
	_, err = qs_t.All(&tags)
	if err != nil {
		logs.Error(err)
	}
	u.Data["posts"] = posts
	u.Data["tags"] = tags
	u.Layout = "layout/base.html"
	u.TplName = "home.html"
}

// 个人中心
func (u *UserController) Profile() {
	//读取flash
	web.ReadFromRequest(&u.Controller)
	flash := web.NewFlash()
	//读取session，看用户是否登录过
	name := u.GetSession("user_name")
	if name == nil {
		flash.Error("请先登录！")
		flash.Store(&u.Controller)
		u.Redirect(web.URLFor("UserController.Home"), 302)
		return
	}
	u.Data["name"] = name.(string)
	var posts []*models.Post
	var tags []*models.Tag
	var comments []*models.Comment
	qs_p := database.Handler.QueryTable("post")
	qs_u := database.Handler.QueryTable("user")
	qs_t := database.Handler.QueryTable("tag")
	qs_c := database.Handler.QueryTable("comment")
	//查询自己的帖子
	user_id := u.GetSession("user_id").(int)
	_, err := qs_p.Filter("Author__Id", user_id).OrderBy("-created").All(&posts)
	if err != nil {
		logs.Error(err)
	}
	user := new(models.User)
	err = qs_u.Filter("id", user_id).One(user)
	if err != nil {
		logs.Error(err)
	}
	qs_c.Filter("Author__Id", user_id).All(&comments)
	for _, comment := range comments {
		tmp := new(models.Post)
		qs_p.Filter("Id", comment.Belong.Id).One(tmp)
		comment.Belong = tmp
	}
	u.Data["posts"] = posts
	u.Data["user"] = user
	qs_t.All(&tags)
	u.Data["tags"] = tags
	u.Data["comments"] = comments
	u.Layout = "layout/base.html"
	u.TplName = "profile.html"
}

// 修改名字
func (u *UserController) Rename() {
	//读取flash
	web.ReadFromRequest(&u.Controller)
	flash := web.NewFlash()
	//读取session，看用户是否登录过
	name := u.GetSession("user_name")
	id := u.GetSession("user_id")
	if name == nil {
		flash.Error("请先登录！")
		flash.Store(&u.Controller)
		u.Redirect(web.URLFor("UserController.Home"), 302)
		return
	}
	//获取新名字
	new_name := u.GetString("name")
	user := new(models.User)
	qs := database.Handler.QueryTable(user)
	//先找到user
	err := qs.Filter("id", id.(int)).One(user)
	if err != nil {
		logs.Error(err)
	}
	//修改名字
	user.Name = new_name
	//数据库更新
	_, err = database.Handler.Update(user, "Name")
	if err != nil {
		logs.Error(err)
		flash.Error("数据库更新失败！")
	} else {
		flash.Success("昵称修改成功。")
		//将session里的改过来
		u.SetSession("user_name", new_name)
	}
	flash.Store(&u.Controller)
	u.Redirect(web.URLFor("UserController.Profile"), 302)
}

// 修改密码
func (u *UserController) ChangePassword() {
	//读取flash
	web.ReadFromRequest(&u.Controller)
	flash := web.NewFlash()
	//读取session，看用户是否登录过
	name := u.GetSession("user_name")
	id := u.GetSession("user_id")
	if name == nil {
		flash.Error("请先登录！")
		flash.Store(&u.Controller)
		u.Redirect(web.URLFor("UserController.Home"), 302)
		return
	}
	user := new(models.User)
	qs := database.Handler.QueryTable(user)
	//先找到user
	err := qs.Filter("id", id.(int)).One(user)
	if err != nil {
		logs.Error(err)
	}
	//先看看密码是否一致
	psw_old := u.GetString("psw-old")
	if !user.CheckPasswordHash(psw_old) {
		flash.Error("请先输入正确的旧密码！")
		flash.Store(&u.Controller)
		u.Redirect(web.URLFor("UserController.Profile"), 302)
		return
	}
	//获取旧密码
	psw_new := u.GetString("psw-new")
	user.Password = psw_new
	user.HashPassword()
	//数据库修改
	database.Handler.Update(user, "Password")
	flash.Success("密码修改成功！")
	flash.Store(&u.Controller)
	u.Redirect(web.URLFor("UserController.Profile"), 302)
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
		user.Avatar = "/static/img/avatar/0.png" //还是这招好使
		//保存到数据库
		_, err = database.Handler.Insert(user)
		if err != nil {
			logs.Error(err)
		}
		flash.Success("注册成功！")
		//填入到session
		u.SetSession("user_name", user.Name)
		qs := database.Handler.QueryTable("user")
		tmp := new(models.User)
		qs.Filter("name", user.Name).One(tmp)
		u.SetSession("user_id", tmp.Id)
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

// 上传头像
func (u *UserController) UploadAvatar() {
	//读取flash
	web.ReadFromRequest(&u.Controller)
	flash := web.NewFlash()
	//读取session，看用户是否登录过
	id := u.GetSession("user_id")
	if id == nil {
		flash.Error("请先登录！")
		flash.Store(&u.Controller)
		u.Redirect(web.URLFor("UserController.Home"), 302)
		return
	}
	user := new(models.User)
	qs := database.Handler.QueryTable(user)
	//先找到user
	err := qs.Filter("id", id.(int)).One(user)
	if err != nil {
		logs.Error(err)
	}
	//获取上传文件
	f, h, err := u.GetFile("avatar")
	if err != nil {
		logs.Error(err)
	}
	defer f.Close()
	//获取文件名后缀
	suffix := path.Ext(h.Filename)
	user.Avatar = "static/img/avatar/" + strconv.Itoa(user.Id) + suffix
	logs.Error(user.Avatar)
	//保存
	u.SaveToFile("avatar", user.Avatar)
	//数据库修改
	user.Avatar = "/" + user.Avatar
	database.Handler.Update(user, "Avatar")
	flash.Success("头像上传成功！")
	flash.Store(&u.Controller)
	u.Redirect(web.URLFor("UserController.Profile"), 302)
}
