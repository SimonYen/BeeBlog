package main

import (
	_ "BeeBlog/routers"

	"github.com/beego/beego/v2/server/web"
)

func main() {
	web.Run()
}
