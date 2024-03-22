package entity

type User struct {
	ID             int64
	Account        string `xorm:"'Account'"`
	Name           string `xorm:"'Name'"`
	Password       string `xorm:"'Password'"`
	Level          int64  `xorm:"'Level'"`          // 1普通用户 2管理员
	Status         int64  `xorm:"'Status'"`         // 1正常 2禁用
	AvailableSpace int64  `xorm:"'AvailableSpace'"` // 可用空间大小M
	UsedSpace      int64  `xorm:"'UsedSpace'"`      // 已使用空间M
	UserToken      string `xorm:"'UserToken'"`
	Email          string `xorm:"'Email'"`
	Captcha        string `xorm:"'Captcha'"`
	Createtime     int64  `xorm:"'Createtime'"`
}
