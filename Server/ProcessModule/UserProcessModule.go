package processmodule

import (
	entity "duckweed-server/Server/Entity"
	lang "duckweed-server/Server/Lang"
	lib "duckweed-server/Server/Lib"
	model "duckweed-server/Server/Model"
)

func SignIn(account string, password string) entity.Result {
	lang := lang.Lang()
	_, _, tx, db := model.ConnDB()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
	}

	b, s, userData := model.UserDataAccount(tx, account)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}

	if userData.ID == 0 {
		tx.Rollback()
		res.Message = lang.NoData
		return res
	}

	if userData.Status == 2 {
		tx.Rollback()
		res.Message = lang.AccountDisabled
		return res
	}

	if lib.MD5(lib.MD5(lib.IntToString(userData.Createtime)+password+lib.IntToString(userData.Createtime))) != userData.Password {
		tx.Rollback()
		res.Message = lang.IncorrectPassword
		return res
	}

	// 写入日志
	b, s = lib.WriteLog(account, account+" login")
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}

	// 检查空间文件夹
	userDir := "../Space/" + account
	if !lib.FileExist(userDir) {
		b, s := lib.DirMake(userDir)
		if !b {
			tx.Rollback()
			res.Message = s
			return res
		}
	}

	// 写入Token
	token := lib.MD5(lib.MD5(lib.TimeNowStr() + userData.Password + lib.TimeNowStr()))
	userData.UserToken = token
	b, s, _ = model.UserUpdate(tx, userData)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}

	tx.Commit()
	res.State = true
	res.Data = userData.UserToken

	db.Close()
	return res
}

func SignOut(userToken string) entity.Result {
	lang := lang.Lang()
	_, _, tx, db := model.ConnDB()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
	}

	if userToken == "" {
		tx.Rollback()
		res.Message = lang.Typo
		return res
	}

	b, s, userData := model.UserDataToken(tx, userToken)
	if userData.ID == 0 {
		tx.Rollback()
		res.Message = lang.NoData
		return res
	}
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}

	userData.UserToken = ""
	b, s, _ = model.UserUpdate(tx, userData)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}

	// 写入日志
	b, s = lib.WriteLog(userData.Account, userData.Account+" logout")
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}

	tx.Commit()
	res.State = true

	db.Close()
	return res
}

func UserList(userToken string, page string, pageSize string, order string, account string, name string, level string, status string) entity.ResultList {
	lang := lang.Lang()
	res := entity.ResultList{
		State:     false,
		Code:      200,
		Message:   "",
		Page:      0,
		PageSize:  0,
		TotalPage: 0,
		Data:      nil,
	}

	permissions, _ := CheckLevel(userToken)
	if permissions != 2 {
		res.Message = lang.NoPermission
		return res
	}

	_, _, tx, db := model.ConnDB()
	_, _, pageInt := lib.StringToInt(page)
	_, _, pageSizeInt := lib.StringToInt(pageSize)
	_, _, orderInt := lib.StringToInt(order)
	_, _, levelInt := lib.StringToInt(level)
	_, _, statusInt := lib.StringToInt(status)
	p, ps, t, list := model.UserList(tx, pageInt, pageSizeInt, orderInt, account, name, levelInt, statusInt)

	res.State = true
	res.Page = p
	res.PageSize = ps
	res.TotalPage = t
	res.Data = list

	tx.Commit()
	db.Close()
	return res
}

func Users(userToken string, order string, account string, name string, level string, status string) entity.Result {
	lang := lang.Lang()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
		Data:    nil,
	}

	permissions, _ := CheckLevel(userToken)
	if permissions != 2 {
		res.Message = lang.NoPermission
		return res
	}

	_, _, tx, db := model.ConnDB()
	_, _, orderInt := lib.StringToInt(order)
	_, _, levelInt := lib.StringToInt(level)
	_, _, statusInt := lib.StringToInt(status)
	list := model.Users(tx, orderInt, account, name, levelInt, statusInt)

	res.State = true
	res.Code = 200
	res.Message = ""
	res.Data = list

	tx.Commit()
	db.Close()
	return res
}

func UserGet(userToken string, id string) entity.Result {
	lang := lang.Lang()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
		Data:    nil,
	}

	permissions, _ := CheckLevel(userToken)
	if permissions != 2 {
		res.Message = lang.NoPermission
		return res
	}

	if id == "" {
		res.Message = lang.Typo
		return res
	}

	_, _, tx, db := model.ConnDB()
	b, s, r := model.UserData(tx, id)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}

	tx.Commit()
	res.Data = r

	db.Close()
	return res
}

func UserAction(userToken string, account string, name string, password string, level string, availableSpace string, id string) entity.Result {
	lang := lang.Lang()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
		Data:    nil,
	}

	permissions, adminAccount := CheckLevel(userToken)
	if permissions != 2 {
		res.Message = lang.NoPermission
		return res
	}

	_, _, levelInt := lib.StringToInt(level)
	_, _, availableSpaceInt := lib.StringToInt(availableSpace)

	if account == "" {
		res.Message = lang.IncorrectAccount
		return res
	}
	if len(account) < 6 {
		res.Message = lang.AccountLengthIsNotEnough
		return res
	}
	if !lib.RegEnNum(account) {
		res.Message = lang.AccountFormatError
		return res
	}
	if name == "" {
		res.Message = lang.IncorrectName
		return res
	}
	if !lib.RegAll(name) {
		res.Message = lang.NicknameFormatError
		return res
	}
	if password == "" {
		res.Message = lang.IncorrectPassword
		return res
	}
	if len(password) < 6 {
		res.Message = lang.PasswordLengthIsNotEnough
		return res
	}
	if !lib.RegEnNum(password) {
		res.Message = lang.PasswordFormatError
		return res
	}
	if level == "" {
		res.Message = lang.IncorrectLevel
		return res
	}
	if levelInt != 1 && levelInt != 2 {
		res.Message = lang.IncorrectLevel
		return res
	}
	if availableSpace == "" {
		res.Message = lang.TheFreeSpaceSizeIsSetIncorrectly
		return res
	}
	if availableSpaceInt < 0 {
		res.Message = lang.TheFreeSpaceSizeIsSetIncorrectly
		return res
	}

	user := entity.UserEntity{}
	user.Account = account
	user.Name = name
	user.Level = levelInt
	user.AvailableSpace = availableSpaceInt
	user.Createtime = int(lib.TimeStamp())

	_, _, tx, db := model.ConnDB()

	var b bool
	var s string
	var r int
	_, _, idInt := lib.StringToInt(id)
	if idInt > 0 {
		if user.ID == 1 {
			user.Level = 2
			user.Status = 1
		}
		if user.Password != password {
			user.Password = lib.MD5(lib.MD5(lib.IntToString(user.Createtime) + password + lib.IntToString(user.Createtime)))
		}

		user.ID = idInt
		b, s, r = model.UserUpdate(tx, user)
		if r == 0 {
			tx.Rollback()
			res.Message = lang.OperationFailed
			return res
		}
	} else {
		// 检查重复账号
		_, _, checkAccount := model.UserDataAccount(tx, account)
		if checkAccount.ID > 0 {
			tx.Rollback()
			res.Message = lang.AccountAlreadyExists
			return res
		}

		user.Password = password
		b, s, r = model.UserAdd(tx, user)
		if r == 0 {
			tx.Rollback()
			res.Message = lang.OperationFailed
			return res
		}
	}

	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}

	if idInt > 0 {
		lib.WriteLog(adminAccount, adminAccount+" update user data account: "+account)
	} else {
		lib.WriteLog(adminAccount, adminAccount+" new user account: "+account)
	}

	tx.Commit()
	res.State = true
	res.Data = r

	db.Close()
	return res
}

func UserDel(userToken string, id string) entity.Result {
	lang := lang.Lang()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
		Data:    nil,
	}

	permissions, adminAccount := CheckLevel(userToken)
	if permissions != 2 {
		res.Message = lang.NoPermission
		return res
	}

	if id == "" {
		res.Message = lang.Typo
		return res
	}

	_, _, tx, db := model.ConnDB()

	b, s, userData := model.UserData(tx, id)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}
	if userData.ID == 0 {
		tx.Rollback()
		res.Message = lang.NoData
		return res
	}

	// 删除用户文件夹数据
	b, s, _ = model.DirDelUser(tx, id)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}

	// 删除用户文件数据
	b, s, _ = model.FileDelUser(tx, id)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}

	// 删除用户数据
	b, s, _ = model.UserDel(tx, id)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}

	// 删除用户磁盘目录
	userDir := "../Space/" + userData.Account
	b, s = lib.DirDel(userDir)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}

	lib.WriteLog(adminAccount, adminAccount+" delete user data account: "+userData.Account)

	tx.Commit()
	res.State = true

	db.Close()
	return res
}

func SignUp(account string, name string, password string) entity.Result {
	lang := lang.Lang()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
		Data:    nil,
	}

	if account == "" {
		res.Message = lang.IncorrectAccount
		return res
	}
	if len(account) < 6 {
		res.Message = lang.AccountLengthIsNotEnough
		return res
	}
	if !lib.RegEnNum(account) {
		res.Message = lang.AccountFormatError
		return res
	}
	if name == "" {
		res.Message = lang.IncorrectName
		return res
	}
	if !lib.RegAll(name) {
		res.Message = lang.NicknameFormatError
		return res
	}
	if password == "" {
		res.Message = lang.IncorrectPassword
		return res
	}
	if len(password) < 6 {
		res.Message = lang.PasswordLengthIsNotEnough
		return res
	}
	if !lib.RegEnNum(password) {
		res.Message = lang.PasswordFormatError
		return res
	}

	_, _, iss := lib.StringToInt(lib.CheckConf().InitialSpaceSize)

	user := entity.UserEntity{}
	user.Account = account
	user.Name = name
	user.Password = password
	user.Level = 1
	user.AvailableSpace = iss
	user.Createtime = int(lib.TimeStamp())

	_, _, tx, db := model.ConnDB()

	_, _, checkAccount := model.UserDataAccount(tx, account)
	if checkAccount.ID > 0 {
		tx.Rollback()
		res.Message = lang.AccountAlreadyExists
		return res
	}

	b, s, r := model.UserAdd(tx, user)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}
	if r == 0 {
		tx.Rollback()
		res.Message = lang.OperationFailed
		return res
	}

	lib.WriteLog(account, "user registration account: "+account)

	tx.Commit()
	res.State = true
	res.Data = r

	db.Close()
	return res
}
