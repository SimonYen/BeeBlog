package main

import (
	_ "BeeBlog/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}

