package models

import (
	"time"

	"github.com/beego/beego/v2/adapter/orm"
)

//注意，不要在数据库中存储原生密码，必须要经过加密
type User struct {
	Id       int
	Name     string
	Email    string
	Password string
	Created  time.Time `orm:"auto_now_add;type(datetime);description(用户创建时间)"`
	Avatar   string    `orm:"size(2048);description(用户头像保持的目录地址)"`
}

func init() {
	orm.RegisterModel(new(User))
}
