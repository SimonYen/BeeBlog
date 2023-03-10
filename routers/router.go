package routers

import (
	"BeeBlog/controllers"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.Router("/", &controllers.UserController{}, "get:Home")
	web.Router("/register", &controllers.UserController{}, "post:Register")
	web.Router("/login", &controllers.UserController{}, "post:Login")
	web.Router("/logout", &controllers.UserController{}, "get:Logout")
	web.Router("/post/save", &controllers.PostController{}, "post:Save")
	web.Router("/post/:id:int", &controllers.PostController{}, "get:Detail")
	web.Router("/profile", &controllers.UserController{}, "get:Profile")
	web.Router("/user/rename", &controllers.UserController{}, "post:Rename")
	web.Router("/user/change_psw", &controllers.UserController{}, "post:ChangePassword")
	web.Router("/upload", &controllers.UserController{}, "post:UploadAvatar")
	web.Router("/post/delete/:id:int", &controllers.PostController{}, "get:Delete")
	web.Router("/tag", &controllers.TagController{}, "post:ChangeTag")
	web.Router("/tag/:id:int", &controllers.TagController{}, "get:TaggedHome")
	web.Router("/comment/add/:id:int", &controllers.CommentController{}, "post:Add")
	web.Router("/comment/delete/:id:int", &controllers.CommentController{}, "get:Delete")
	web.Router("/search", &controllers.PostController{}, "get:Search")
}
