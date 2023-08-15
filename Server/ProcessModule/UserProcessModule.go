package processmodule

import (
	entity "duckweed-server/Server/Entity"
	lang "duckweed-server/Server/Lang"
	lib "duckweed-server/Server/Lib"
	model "duckweed-server/Server/Model"
)

/*
func UserGet(id string) entity.Result {
	_, _, db := model.ConnDB()
	res := entity.Result{
		State:   true,
		Code:    200,
		Message: "",
	}
	b, s, r := model.UserData(db, id)
	if !b {
		res.State = false
		res.Message = s
	} else {
		res.Data = r
	}
	db.Close()
	return res
}

func Users(order string, account string, name string, level string, status string) entity.Result {
	_, _, db := model.ConnDB()
	_, _, orderInt := lib.StringToInt(order)
	_, _, levelInt := lib.StringToInt(level)
	_, _, statusInt := lib.StringToInt(status)
	res := entity.Result{
		State:   true,
		Code:    200,
		Message: "",
		Data:    model.Users(db, orderInt, account, name, levelInt, statusInt),
	}
	db.Close()
	return res
}

func UserList(page string, pageSize string, order string, account string, name string, level string, status string) entity.ResultList {
	_, _, db := model.ConnDB()
	_, _, pageInt := lib.StringToInt(page)
	_, _, pageSizeInt := lib.StringToInt(pageSize)
	_, _, orderInt := lib.StringToInt(order)
	_, _, levelInt := lib.StringToInt(level)
	_, _, statusInt := lib.StringToInt(status)
	p, ps, t, list := model.UserList(db, pageInt, pageSizeInt, orderInt, account, name, levelInt, statusInt)
	res := entity.ResultList{
		State:     true,
		Code:      200,
		Message:   "",
		Page:      p,
		PageSize:  ps,
		TotalPage: t,
		Data:      list,
	}
	db.Close()
	return res
}

func UserCheck(account string, name string, password string, level string, availableSpace string, id string) entity.Result {
	_, _, db := model.ConnDB()
	res := entity.Result{
		State:   true,
		Code:    200,
		Message: "",
	}
	_, _, levelInt := lib.StringToInt(level)
	_, _, availableSpaceInt := lib.StringToInt(availableSpace)
	user := entity.UserEntity{
		Account:        account,
		Name:           name,
		Password:       password,
		Level:          levelInt,
		Status:         1,
		AvailableSpace: availableSpaceInt,
	}

	var b bool
	var s string
	var r int
	_, _, idInt := lib.StringToInt(id)
	if idInt > 0 {
		b, s, r = model.UserUpdate(db, id, user)
	} else {
		b, s, r = model.UserAdd(db, user)
	}

	if !b {
		res.State = false
		res.Message = s
	} else {
		res.Data = r
	}

	res.Data = r
	db.Close()
	return res
}

func UserDel(id string) entity.Result {
	_, _, db := model.ConnDB()
	res := entity.Result{
		State:   true,
		Code:    200,
		Message: "",
	}
	b, s, r := model.UserDel(db, id)
	if !b {
		res.State = false
		res.Message = s
	} else {
		res.Data = r
	}
	db.Close()
	return res
}
*/

func UserLogin(account string, password string) entity.Result {
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
		res.Message = lang["NoData"]
		return res
	}

	if lib.MD5(lib.MD5(lib.IntToString(userData.Createtime)+password+lib.IntToString(userData.Createtime))) != userData.Password {
		tx.Rollback()
		res.Message = lang["IncorrectPassword"]
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
