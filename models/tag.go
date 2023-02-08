package models

type Tag struct {
	Id   int
	Name string `orm:"unique"`

	//简化了一下，一个文章只能有一个分类，所以是一对多的关系
	Posts []*Post `orm:"reverse(many)"`
}
