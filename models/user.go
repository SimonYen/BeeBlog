package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// 注意，不要在数据库中存储原生密码，必须要经过加密
type User struct {
	Id       int       `form:"-"`
	Name     string    `form:"name"`
	Email    string    `form:"email" orm:"unique"`
	Password string    `form:"psw"`
	Created  time.Time `orm:"auto_now_add;type(datetime);description(用户创建时间)" form:"-"`
	Avatar   string    `orm:"size(2048);description(用户头像保持的目录地址);default(/static/img/avatar/0.png)" form:"-"`

	Posts []*Post `orm:"reverse(many)"`
}

// 密码哈希加密
func (receiver *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(receiver.Password), 14)
	receiver.Password = string(bytes)
	return err
}

// 密码比对
func (receiver *User) CheckPasswordHash(passsword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(receiver.Password), []byte(passsword))
	return err == nil
}
