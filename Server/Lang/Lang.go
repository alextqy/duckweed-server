package lang

import lib "duckweed-server/Server/Lib"

type language struct {
	NoData                           string
	IncorrectPassword                string
	PasswordLengthIsNotEnough        string
	PasswordFormatError              string
	Typo                             string
	NoPermission                     string
	AccountDisabled                  string
	IncorrectAccount                 string
	AccountLengthIsNotEnough         string
	AccountAlreadyExists             string
	AccountFormatError               string
	IncorrectName                    string
	NicknameFormatError              string
	IncorrectLevel                   string
	TheFreeSpaceSizeIsSetIncorrectly string
	OperationFailed                  string
	MalformedContent                 string
}

func Lang() language {
	l := language{}
	if lib.CheckConf().Lang == "zh" {
		l.NoData = "无数据"
		l.IncorrectPassword = "密码错误"
		l.PasswordLengthIsNotEnough = "密码长度不够"
		l.PasswordFormatError = "密码格式错误"
		l.Typo = "输入错误"
		l.NoPermission = "无权限"
		l.AccountDisabled = "账号已禁用"
		l.IncorrectAccount = "账号错误"
		l.AccountLengthIsNotEnough = "账号长度不够"
		l.AccountAlreadyExists = "账号已存在"
		l.AccountFormatError = "账号格式错误"
		l.IncorrectName = "昵称错误"
		l.NicknameFormatError = "昵称格式错误"
		l.IncorrectLevel = "等级错误"
		l.TheFreeSpaceSizeIsSetIncorrectly = "可用空间设置错误"
		l.OperationFailed = "操作失败"
		l.MalformedContent = "内容格式错误"
	} else if lib.CheckConf().Lang == "en" {
		l.NoData = "no data"
		l.IncorrectPassword = "incorrect password"
		l.PasswordLengthIsNotEnough = "password length is not enough"
		l.PasswordFormatError = "password format error"
		l.Typo = "typo"
		l.NoPermission = "no permission"
		l.AccountDisabled = "account disabled"
		l.IncorrectAccount = "incorrect account"
		l.AccountLengthIsNotEnough = "account length is not enough"
		l.AccountAlreadyExists = "account already exists"
		l.AccountFormatError = "account format error"
		l.IncorrectName = "incorrect name"
		l.NicknameFormatError = "nick name format error"
		l.IncorrectLevel = "incorrect level"
		l.TheFreeSpaceSizeIsSetIncorrectly = "the free space size is set incorrectly"
		l.OperationFailed = "operation failed"
		l.MalformedContent = "malformed content"
	} else {
		l.NoData = ""
		l.IncorrectPassword = ""
		l.PasswordLengthIsNotEnough = ""
		l.PasswordFormatError = ""
		l.Typo = ""
		l.NoPermission = ""
		l.AccountDisabled = ""
		l.IncorrectAccount = ""
		l.AccountLengthIsNotEnough = ""
		l.AccountAlreadyExists = ""
		l.AccountFormatError = ""
		l.IncorrectName = ""
		l.NicknameFormatError = ""
		l.IncorrectLevel = ""
		l.TheFreeSpaceSizeIsSetIncorrectly = ""
		l.OperationFailed = ""
		l.MalformedContent = ""
	}
	return l
}
