package entity

type FileEntity struct {
	ID            int
	FileName      string
	FileType      string
	FileSize      string
	StoragePath   string
	MD5           string
	UserID        int
	DirID         int // -1 为回收站
	Createtime    int
	Status        int    // 1 上传中 2 完成 3 操作失败
	OutreachID    string // 唯一编码
	SourceAddress string // 上传地址
}
