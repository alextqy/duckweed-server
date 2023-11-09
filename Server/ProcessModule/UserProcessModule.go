package processmodule

import (
	entity "duckweed-server/Server/Entity"
	lang "duckweed-server/Server/Lang"
	lib "duckweed-server/Server/Lib"
	model "duckweed-server/Server/Model"
)

func SignIn(account, password string) entity.Result {
	lang := lang.Lang()
	_, _, tx, db := model.ConnDB()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
	}

	if account == "" {
		res.Message = lang.Typo
		return res
	}

	if password == "" {
		res.Message = lang.Typo
		return res
	}

	b, s, userData := model.UserDataAccount(tx, account)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}

	if userData.ID == 0 {
		tx.Rollback()
		res.Message = lang.AccountDoesNotExist
		return res
	}

	if userData.Status == 2 {
		tx.Rollback()
		res.Message = lang.AccountDisabled
		return res
	}

	if UserPWD(userData, password) != userData.Password {
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
	token := UserToken(userData)
	userData.UserToken = token
	userData.Captcha = ""
	b, s, _ = model.UserUpdate(tx, userData)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}

	res.State = true
	res.Message = lib.IntToString(userData.Level)
	res.Data = userData.UserToken

	tx.Commit()
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

	res.State = true

	tx.Commit()
	db.Close()
	return res
}

// 管理员操作 =============================================================================================================================================

func UserList(userToken, page, pageSize, order, account, name, level, status string) entity.ResultList {
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

	permissions, _, _ := CheckLevel(userToken)
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

func Users(userToken, order, account, name, level, status string) entity.Result {
	lang := lang.Lang()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
		Data:    nil,
	}

	permissions, _, _ := CheckLevel(userToken)
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

func UserGet(userToken, id string) entity.Result {
	lang := lang.Lang()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
		Data:    nil,
	}

	permissions, _, _ := CheckLevel(userToken)
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

	res.State = true
	res.Data = r

	tx.Commit()
	db.Close()
	return res
}

func SetAvailableSpace(userToken, id, availableSpace string) entity.Result {
	lang := lang.Lang()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
		Data:    nil,
	}

	if id == "" {
		res.Message = lang.Typo
		return res
	}
	if availableSpace == "" {
		res.Message = lang.Typo
		return res
	}
	_, _, availableSpaceInt := lib.StringToInt(availableSpace)
	if availableSpaceInt <= 0 {
		res.Message = lang.Typo
		return res
	}

	permissions, adminAccount, _ := CheckLevel(userToken)
	if permissions != 2 {
		res.Message = lang.NoPermission
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

	userData.AvailableSpace = availableSpaceInt
	b, s, r := model.UserUpdate(tx, userData)
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

	lib.WriteLog(adminAccount, "set user available space account: "+userData.Account+" "+availableSpace+"M")

	res.State = true

	tx.Commit()
	db.Close()
	return res
}

func SetRootAccount(userToken, id string) entity.Result {
	lang := lang.Lang()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
		Data:    nil,
	}

	if id == "" {
		res.Message = lang.Typo
		return res
	}

	permissions, adminAccount, rootID := CheckLevel(userToken)
	if permissions != 2 {
		res.Message = lang.NoPermission
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
	_, _, idInt := lib.StringToInt(id)
	if rootID == idInt {
		tx.Rollback()
		res.Message = lang.OperationFailed
		return res
	}

	if id == "1" {
		userData.Level = 2
	} else {
		if userData.Level == 1 {
			userData.Level = 2
		} else if userData.Level == 2 {
			userData.Level = 1
		} else {
			userData.Level = 1
		}
	}

	b, s, r := model.UserUpdate(tx, userData)
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

	if userData.Level == 2 {
		lib.WriteLog(adminAccount, "set user root account: "+userData.Account)
	} else {
		lib.WriteLog(adminAccount, "unset user root account: "+userData.Account)
	}

	res.State = true

	tx.Commit()
	db.Close()
	return res
}

func DisableUser(userToken, id string) entity.Result {
	lang := lang.Lang()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
		Data:    nil,
	}

	if id == "" {
		res.Message = lang.Typo
		return res
	}

	permissions, adminAccount, rootID := CheckLevel(userToken)
	if permissions != 2 {
		res.Message = lang.NoPermission
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
	_, _, idInt := lib.StringToInt(id)
	if rootID == idInt {
		tx.Rollback()
		res.Message = lang.OperationFailed
		return res
	}

	if id == "1" {
		userData.Status = 1
	} else {
		if userData.Status == 1 {
			userData.Status = 2
		} else if userData.Status == 2 {
			userData.Status = 1
		} else {
			userData.Status = 1
		}
	}

	b, s, r := model.UserUpdate(tx, userData)
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

	if userData.Status == 1 {
		lib.WriteLog(adminAccount, "enable user data account: "+userData.Account)
	} else {
		lib.WriteLog(adminAccount, "disable user data account: "+userData.Account)
	}

	res.State = true

	tx.Commit()
	db.Close()
	return res
}

func UserDel(userToken, id string) entity.Result {
	lang := lang.Lang()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
		Data:    nil,
	}

	permissions, adminAccount, _ := CheckLevel(userToken)
	if permissions != 2 {
		res.Message = lang.NoPermission
		return res
	}

	if id == "" {
		res.Message = lang.Typo
		return res
	}
	if id == "1" {
		res.Message = lang.OperationFailed
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

	// 删除用户日志文件
	b, s, d, _ := lib.DirTraverse("../Log")
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}
	if len(d) > 0 {
		for i := 0; i < len(d); i++ {
			f := d[i] + "/" + userData.Account + ".log"
			lib.FileRemove(f)
		}
	}

	lib.WriteLog(adminAccount, "delete user data account: "+userData.Account)

	res.State = true

	tx.Commit()
	db.Close()
	return res
}

// 用户操作 =============================================================================================================================================

func SignUp(account, name, password, email, captcha string) entity.Result {
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
	if len(account) < 4 {
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
	if email == "" {
		res.Message = lang.EmailError
		return res
	}
	if !lib.RegEmail(email) {
		res.Message = lang.EmailFormatError
		return res
	}
	if captcha == "" {
		res.Message = lang.IncorrectCaptcha
		return res
	}
	if !lib.FileExist("../Temp/Captcha/" + captcha) {
		res.Message = lang.IncorrectCaptcha
		return res
	}
	b, s := lib.FileRead("../Temp/Captcha/" + captcha)
	if !b {
		res.Message = s
		return res
	}
	if s != email {
		res.Message = lang.EmailError
		return res
	}

	lib.FileRemove("../Temp/Captcha/" + captcha)

	_, _, iss := lib.StringToInt(lib.CheckConf().InitialSpaceSize)

	user := entity.UserEntity{}
	user.Account = account
	user.Name = name
	user.Password = password
	user.Level = 1
	user.AvailableSpace = iss
	user.Createtime = int(lib.TimeStamp())
	user.Email = email

	_, _, tx, db := model.ConnDB()

	b, s, checkEmail := model.UserDataEmail(tx, email)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}
	if checkEmail.ID > 0 {
		tx.Rollback()
		res.Message = lang.EmailAlreadyInUse
		return res
	}

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

	res.State = true
	res.Data = r

	tx.Commit()
	db.Close()
	return res
}

func CheckPersonalData(userToken string) entity.Result {
	lang := lang.Lang()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
		Data:    nil,
	}
	userData := CheckToken(userToken)
	if userData.ID == 0 {
		res.Message = lang.ReLoginRequired
		return res
	}
	res.State = true
	res.Data = userData
	return res
}

func ModifyPersonalData(userToken, name, password, email, captcha string) entity.Result {
	lang := lang.Lang()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
		Data:    nil,
	}

	if name == "" {
		res.Message = lang.IncorrectName
		return res
	}
	if !lib.RegAll(name) {
		res.Message = lang.NicknameFormatError
		return res
	}
	if email == "" {
		res.Message = lang.EmailError
		return res
	}

	userData := CheckToken(userToken)
	if userData.ID == 0 {
		res.Message = lang.ReLoginRequired
		return res
	}

	userData.Name = name
	if password != "" && password != userData.Password {
		if len(password) < 6 {
			res.Message = lang.PasswordLengthIsNotEnough
			return res
		}
		if !lib.RegEnNum(password) {
			res.Message = lang.PasswordFormatError
			return res
		}
		userData.Password = UserPWD(userData, password)
	}

	_, _, tx, db := model.ConnDB()
	if userData.Email != email {
		if !lib.RegEmail(email) {
			res.Message = lang.EmailFormatError
			return res
		}
		if captcha == "" {
			tx.Rollback()
			res.Message = lang.IncorrectCaptcha
			return res
		}
		if userData.Captcha != captcha {
			tx.Rollback()
			res.Message = lang.IncorrectCaptcha
			return res
		}
		b, s, checkEmail := model.UserDataEmail(tx, email)
		if !b {
			tx.Rollback()
			res.Message = s
			return res
		}
		if checkEmail.ID > 0 {
			tx.Rollback()
			res.Message = lang.EmailAlreadyInUse
			return res
		}
		userData.Email = email
	}

	b, s, r := model.UserUpdate(tx, userData)
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

	lib.WriteLog(userData.Account, "user modify personal information account: "+userData.Account)

	res.State = true
	res.Data = r

	tx.Commit()
	db.Close()
	return res
}

func ResetPassword(newPassword, captcha string) entity.Result {
	lang := lang.Lang()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
		Data:    nil,
	}

	if newPassword == "" {
		res.Message = lang.IncorrectPassword
		return res
	}
	if len(newPassword) < 6 {
		res.Message = lang.PasswordLengthIsNotEnough
		return res
	}
	if !lib.RegEnNum(newPassword) {
		res.Message = lang.PasswordFormatError
		return res
	}
	if captcha == "" {
		res.Message = lang.Typo
		return res
	}

	_, _, tx, db := model.ConnDB()

	b, s, userData := model.UserDataCaptcha(tx, captcha)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}
	if userData.ID == 0 {
		tx.Rollback()
		res.Message = lang.IncorrectCaptcha
		return res
	}

	userData.Password = UserPWD(userData, newPassword)
	userData.Captcha = ""
	b, s, r := model.UserUpdate(tx, userData)
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

	res.State = true

	tx.Commit()
	db.Close()
	return res
}

func SendEmail(email string, sendType string) entity.Result {
	lang := lang.Lang()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
		Data:    nil,
	}

	if email == "" {
		res.Message = lang.Typo
		return res
	}
	if !lib.RegEmail(email) {
		res.Message = lang.EmailFormatError
		return res
	}
	if sendType == "" {
		res.Message = lang.IncorrectSendingType
		return res
	}
	if !lib.RegNum(sendType) {
		res.Message = lang.IncorrectSendingType
		return res
	}

	_, _, tx, db := model.ConnDB()

	captcha := lib.RandStr(5) // 验证码

	_, _, sendTypeNum := lib.StringToInt(sendType)
	if sendTypeNum <= 0 {
		res.Message = lang.IncorrectSendingType
		return res
	}

	subject := ""
	b, s, userData := model.UserDataEmail(tx, email)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}
	if sendTypeNum == 1 {
		subject = "User registration"

		captchaFile := "../Temp/Captcha/" + captcha
		b, s := lib.FileMake(captchaFile)
		if !b {
			res.Message = s
			return res
		}
		b, s = lib.FileWrite(captchaFile, email)
		if !b {
			res.Message = s
			return res
		}
		if userData.ID > 0 {
			tx.Rollback()
			res.Message = lang.EmailAlreadyInUse
			return res
		}
	} else if sendTypeNum == 2 {
		subject = "Reset Password"

		if userData.ID == 0 {
			tx.Rollback()
			res.Message = lang.EmailAddressDoesNotExist
			return res
		}

		userData.Captcha = captcha
		b, s, r := model.UserUpdate(tx, userData)
		if !b {
			tx.Rollback()
			res.Message = s
			return res
		}
		if r == 0 {
			tx.Rollback()
			res.Message = s
			return res
		}
	} else if sendTypeNum == 3 {
		subject = "Modify email"

		if userData.ID > 0 {
			tx.Rollback()
			res.Message = lang.EmailAlreadyInUse
			return res
		}

		userData.Captcha = captcha
		b, s, r := model.UserUpdate(tx, userData)
		if !b {
			tx.Rollback()
			res.Message = s
			return res
		}
		if r == 0 {
			tx.Rollback()
			res.Message = s
			return res
		}
	} else {
		tx.Rollback()
		res.Message = lang.OperationFailed
		return res
	}

	b, s = lib.SendEmail("tqyalex@qq.com", "qfjhhammjflgbjcc", "Duckweed Server", "smtp.qq.com:587", email, subject, captcha) // 发送邮件
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}

	res.State = true

	tx.Commit()
	db.Close()
	return res
}
