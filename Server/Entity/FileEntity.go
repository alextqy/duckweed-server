package entity

type File struct {
	Id            int64
	FileName      string `xorm:"'FileName'"`
	FileType      string `xorm:"'FileType'"`
	FileSize      string `xorm:"'FileSize'"`
	StoragePath   string `xorm:"'StoragePath'"`
	MD5           string `xorm:"'MD5'"`
	UserID        int64  `xorm:"'UserID'"`
	DirID         int64  `xorm:"'DirID'"`         // -1 为回收站
	Status        int64  `xorm:"'Status'"`        // 1 上传中 2 完成 3 操作失败
	OutreachID    string `xorm:"'OutreachID'"`    // 唯一编码
	SourceAddress string `xorm:"'SourceAddress'"` // 上传地址
	Createtime    int64  `xorm:"'Createtime'"`
}
