package database

import (
	"BeeBlog/models"
	"fmt"

	"github.com/beego/beego/v2/adapter/orm"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

// 创建ORM实例，用于整个APP里复用
var Handler orm.Ormer

// 方便得到类别
var id2name map[int]string

func init() {
	//准备连接数据库

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", GetDataSource())
	orm.RegisterModel(new(models.User))
	orm.RegisterModel(new(models.Post))
	orm.RegisterModel(new(models.Tag))
	orm.RegisterModel(new(models.Comment))

	//模型迁移
	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		panic(err)
	}
	Handler = orm.NewOrm()
	id2name = make(map[int]string)
	qs := Handler.QueryTable("tag")
	var tags []*models.Tag
	qs.All(&tags)
	//存入到字典中
	for _, tag := range tags {
		id2name[tag.Id] = tag.Name
	}
}

// 打算将session存在mysql中，因此把dataSource引出去比较方便
func GetDataSource() string {

	username, _ := web.AppConfig.String("username")
	password, _ := web.AppConfig.String("password")
	host, _ := web.AppConfig.String("host")
	port, _ := web.AppConfig.Int("port")
	database, _ := web.AppConfig.String("database")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&", username, password, host, port, database)
	dataSource += "loc=Asia%2FShanghai"

	return dataSource
}

func GetTagName(id int) string {
	name, ok := id2name[id]
	if ok {
		return name
	}
	return "未分类"
}
