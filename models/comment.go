package models

import "time"

type Comment struct {
	Id      int       `form:"-"`
	Content string    `orm:"type(text)" form:"content"`
	Created time.Time `orm:"auto_now_add;type(datetime);description(文章创建时间)" form:"-"`

	//所属文章
	Belong *Post `orm:"rel(fk)"`
	Author *User `orm:"rel(fk)"`
}
