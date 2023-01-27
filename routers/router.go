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
	web.Router("/post/create", &controllers.PostController{}, "get:Create")
	web.Router("/post/save", &controllers.PostController{}, "post:Save")
}
