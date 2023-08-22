package entity

type UserEntity struct {
	ID             int
	Account        string
	Name           string
	Password       string
	Level          int // 1普通用户 2管理员
	Status         int // 1正常 2禁用
	AvailableSpace int // 可用空间大小M
	UsedSpace      int // 已使用空间M
	Createtime     int
	UserToken      string
	Email          string
	Captcha        string
}
