package main

import (
	_ "BeeBlog/routers"
	"fmt"

	_ "BeeBlog/models"

	"github.com/beego/beego/v2/adapter/orm"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	//准备连接数据库
	username, _ := web.AppConfig.String("username")
	password, _ := web.AppConfig.String("password")
	host, _ := web.AppConfig.String("host")
	port, _ := web.AppConfig.Int("port")
	database, _ := web.AppConfig.String("database")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4", username, password, host, port, database)

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", dataSource)

	//模型迁移
	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		panic(err)
	}
}

func main() {
	if web.BConfig.RunMode == "dev" {
		orm.Debug = true
	}
	orm.RunCommand()
	web.Run()
}
