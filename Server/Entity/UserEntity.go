package entity

type UserEntity struct {
	ID         int
	Account    string
	Name       string
	Password   string
	Level      int // 1普通用户 2管理员
	Status     int // 1正常 2禁用
	Createtime int
}
