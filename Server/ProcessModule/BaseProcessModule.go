package processmodule

import (
	entity "duckweed-server/Server/Entity"
	model "duckweed-server/Server/Model"
)

func CheckLevel(userToken string) (int, string) {
	_, _, tx, db := model.ConnDB()
	if userToken == "" {
		tx.Rollback()
		return 0, ""
	}
	b, _, userData := model.UserDataToken(tx, userToken)
	if !b {
		tx.Rollback()
		return 0, ""
	}
	if userData.Status == 2 {
		tx.Rollback()
		return 0, ""
	}
	tx.Commit()
	db.Close()
	return userData.Level, userData.Account
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
