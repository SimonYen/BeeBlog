package models

import "time"

type Post struct {
	Id          int
	Title       string
	Description string
	Content     string    `orm:"type(text)"`
	Cover       string    `orm:"description(文章封面图)"`
	ReadNum     int       `orm:"description(阅读数);default(0)" form:"-"`
	StarNum     int       `orm:"description(点赞数);default(0)" form:"-"`
	Created     time.Time `orm:"auto_now_add;type(datetime);description(文章创建时间)" form:"-"`
	Updated     time.Time `orm:"auto_now;type(datetime);description(文章修改时间)" form:"-"`

	Author *User `orm:"rel(fk)"`
}
