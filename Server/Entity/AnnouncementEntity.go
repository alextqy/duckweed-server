package entity

type Announcement struct {
	ID         int64
	Content    string `xorm:"'Content'"`
	Createtime int64  `xorm:"'Createtime'"`
}
