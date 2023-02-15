package models

import (
	"strings"
	"time"
)

type Post struct {
	Id          int       `form:"-"`
	Title       string    `form:"title"`
	Description string    `form:"desc"`
	Content     string    `orm:"type(text)" form:"content"`
	Created     time.Time `orm:"auto_now_add;type(datetime);description(文章创建时间)" form:"-"`
	Updated     time.Time `orm:"auto_now;type(datetime);description(文章修改时间)" form:"-"`

	Author *User `orm:"rel(fk)"`
	Class  *Tag  `orm:"rel(fk)"`

	Comments []*Comment `orm:"reverse(many)"`
}

// 主要是为了将string中的换行符给替换<br>
func (p *Post) ToHTML() string {
	result := strings.ReplaceAll(p.Content, "\r\n", "<br>")
	result = strings.ReplaceAll(result, "\n", "<br>")
	result = strings.ReplaceAll(result, "\r", "<br>")
	return result
}
