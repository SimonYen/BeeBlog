package routers

import (
	"BeeBlog/controllers"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.Router("/", &controllers.HomeController{})
	web.Router("/register", &controllers.RegisterController{})
	web.Router("/login", &controllers.LoginController{})
	web.Router("/logout", &controllers.LogoutController{})
	web.Router("/post/create", &controllers.PostController{}, "get:Create")
	web.Router("/post/save", &controllers.PostController{}, "post:Save")
}
