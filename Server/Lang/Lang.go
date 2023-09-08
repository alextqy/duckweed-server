package lang

import lib "duckweed-server/Server/Lib"

type language struct {
	ReLoginRequired                  string
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
	DirectoryAlreadyExists           string
	DirectoryDoesNotExist            string
	ParentFolderDoesNotExist         string
	WrongFormatOfFolderName          string
	FileNameFormatError              string
	FileAlreadyExists                string
	EmailError                       string
	EmailFormatError                 string
	EmailAlreadyInUse                string
	EmailAddressDoesNotExist         string
	IncorrectCaptcha                 string
	LogDoesNotExist                  string
	IncorrectSendingType             string
}

func Lang() language {
	l := language{}
	if lib.CheckConf().Lang == "zh" {
		l.ReLoginRequired = "需要重新登录"
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
		l.DirectoryAlreadyExists = "文件夹已存在"
		l.DirectoryDoesNotExist = "文件夹不存在"
		l.ParentFolderDoesNotExist = "上级文件夹不存在"
		l.WrongFormatOfFolderName = "文件夹名称格式错误"
		l.FileNameFormatError = "文件名称格式错误"
		l.FileAlreadyExists = "文件已存在"
		l.EmailError = "电子邮件错误"
		l.EmailFormatError = "电子邮件格式错误"
		l.EmailAlreadyInUse = "电子邮件已被使用"
		l.EmailAddressDoesNotExist = "电子邮箱不存在"
		l.IncorrectCaptcha = "验证码不正确"
		l.LogDoesNotExist = "日志不存在"
		l.IncorrectSendingType = "发送类型错误"
	} else if lib.CheckConf().Lang == "en" {
		l.ReLoginRequired = "re-login required"
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
		l.DirectoryAlreadyExists = "directory already exists"
		l.DirectoryDoesNotExist = "directory does not exist"
		l.ParentFolderDoesNotExist = "parent folder does not exist"
		l.WrongFormatOfFolderName = "wrong format of folder name"
		l.FileNameFormatError = "file name format error"
		l.FileAlreadyExists = "file already exists"
		l.EmailError = "email error"
		l.EmailFormatError = "email format error"
		l.EmailAlreadyInUse = "email already in use"
		l.EmailAddressDoesNotExist = "email address does not exist"
		l.IncorrectCaptcha = "incorrect captcha"
		l.LogDoesNotExist = "log does not exist"
		l.IncorrectSendingType = "incorrect sending type"
	} else {
		l.ReLoginRequired = ""
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
		l.DirectoryAlreadyExists = ""
		l.DirectoryDoesNotExist = ""
		l.ParentFolderDoesNotExist = ""
		l.WrongFormatOfFolderName = ""
		l.FileNameFormatError = ""
		l.FileAlreadyExists = ""
		l.EmailError = ""
		l.EmailFormatError = ""
		l.EmailAlreadyInUse = ""
		l.EmailAddressDoesNotExist = ""
		l.IncorrectCaptcha = ""
		l.LogDoesNotExist = ""
		l.IncorrectSendingType = ""
	}
	return l
}
