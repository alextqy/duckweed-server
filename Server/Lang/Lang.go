package lang

import (
	entity "duckweed-server/Server/Entity"
	lib "duckweed-server/Server/Lib"
	"encoding/json"
)

type language struct {
	NoData            string
	IncorrectPassword string
	Typo              string
	NoPermission      string
	AccountDisabled   string
}

func Lang() language {
	var confEntity entity.ConfEntity
	_, byteData := lib.FileRead("./Conf.json")
	json.Unmarshal([]byte(byteData), &confEntity)

	l := language{}
	if confEntity.Lang == "zh" {
		l.NoData = "无数据"
		l.IncorrectPassword = "密码错误"
		l.Typo = "输入有误"
		l.NoPermission = "无权限"
		l.AccountDisabled = "账号已禁用"
	} else if confEntity.Lang == "en" {
		l.NoData = "no data"
		l.IncorrectPassword = "incorrect password"
		l.Typo = "typo"
		l.NoPermission = "no permission"
		l.AccountDisabled = "account disabled"
	} else {
		l.NoData = ""
		l.IncorrectPassword = ""
		l.Typo = ""
		l.NoPermission = ""
		l.AccountDisabled = ""
	}
	return l
}
