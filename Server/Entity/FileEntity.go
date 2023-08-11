package entity

type FileEntity struct {
	ID          int
	FileName    string
	FileType    string
	FileSize    string
	StoragePath string
	MD5         string
	UserID      int
	DirID       int
	Createtime  int
}
