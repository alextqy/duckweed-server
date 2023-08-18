package entity

type DirEntity struct {
	ID         int
	DirName    string
	ParentID   int // -1 为回收站
	UserID     int
	Createtime int
}
