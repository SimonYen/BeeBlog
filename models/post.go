package models

import "time"

type Post struct {
	Id          int       `form:"-"`
	Title       string    `form:"title"`
	Description string    `form:"desc"`
	Content     string    `orm:"type(text)" form:"content"`
	Cover       string    `orm:"description(文章封面图)"`
	Created     time.Time `orm:"auto_now_add;type(datetime);description(文章创建时间)" form:"-"`
	Updated     time.Time `orm:"auto_now;type(datetime);description(文章修改时间)" form:"-"`

	Author *User `orm:"rel(fk)"`
}
