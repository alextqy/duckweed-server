package processmodule

import (
	entity "duckweed-server/Server/Entity"
	lib "duckweed-server/Server/Lib"
	model "duckweed-server/Server/Model"
)

func CheckLevel(userToken string) (int, string, int) {
	_, _, tx, db := model.ConnDB()
	if userToken == "" {
		tx.Rollback()
		return 0, "", 0
	}
	b, _, userData := model.UserDataToken(tx, userToken)
	if !b {
		tx.Rollback()
		return 0, "", 0
	}
	if userData.ID == 0 {
		tx.Rollback()
		return 0, "", 0
	}
	if userData.Status == 2 {
		tx.Rollback()
		return 0, "", 0
	}
	tx.Commit()
	db.Close()
	return userData.Level, userData.Account, userData.ID
}

func CheckToken(userToken string) entity.UserEntity {
	userData := entity.UserEntity{}

	_, _, tx, db := model.ConnDB()
	if userToken == "" {
		tx.Rollback()
		return userData
	}
	b, _, userData := model.UserDataToken(tx, userToken)
	if !b {
		tx.Rollback()
		return userData
	}
	if userData.Status == 2 {
		tx.Rollback()
		return userData
	}

	tx.Commit()
	db.Close()
	return userData
}

func UserPWD(data entity.UserEntity, password string) string {
	return lib.MD5(lib.MD5(lib.IntToString(data.Createtime) + data.Account + password + lib.IntToString(data.Createtime)))
}

func UserToken(data entity.UserEntity) string {
	return lib.MD5(lib.MD5(lib.TimeNowStr() + data.Password + lib.RandStr(4) + lib.TimeNowStr()))
}
