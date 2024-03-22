package entity

type Dir struct {
	ID         int64
	DirName    string `xorm:"'DirName'"`
	ParentID   int64  `xorm:"'ParentID'"` // -1 为回收站
	UserID     int64  `xorm:"'UserID'"`
	Createtime int64  `xorm:"'Createtime'"`
}
