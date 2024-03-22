package processmodule

import (
	"database/sql"
	entity "duckweed-server/Server/Entity"
	lang "duckweed-server/Server/Lang"
	lib "duckweed-server/Server/Lib"
	model "duckweed-server/Server/Model"
	"strings"
)

func Dirs(userToken, order, parentID, dirName string) entity.Result {
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

	_, _, orderInt := lib.StringToInt(order)
	_, _, parentIDInt := lib.StringToInt64(parentID)
	if parentIDInt < 0 {
		res.Message = lang.Typo
		return res
	}

	_, _, tx, db := model.ConnDB()

	res.State = true
	res.Data = model.Dirs(tx, orderInt, dirName, parentIDInt, userData.ID)

	tx.Commit()
	db.Close()
	return res
}

func DirAction(userToken, dirName, parentID, id string) entity.Result {
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

	if dirName == "" {
		res.Message = lang.Typo
		return res
	}
	if !lib.RegAll(dirName) {
		res.Message = lang.WrongFormatOfFolderName
		return res
	}
	if parentID == "" {
		res.Message = lang.Typo
		return res
	}

	_, _, idInt := lib.StringToInt64(id)
	if idInt < 0 {
		res.Message = lang.Typo
		return res
	}

	_, _, parentIDInt := lib.StringToInt64(parentID)
	if parentIDInt < 0 {
		res.Message = lang.Typo
		return res
	}

	_, _, tx, db := model.ConnDB()

	if parentIDInt > 0 {
		if parentIDInt == idInt {
			res.Message = lang.OperationFailed
			return res
		}

		b, s, r := model.DirData(tx, parentID)
		if !b {
			tx.Rollback()
			res.Message = s
			return res
		}
		if r.ID == 0 {
			tx.Rollback()
			res.Message = lang.ParentFolderDoesNotExist
			return res
		}
	}

	b, s, sd := model.DirDataSame(tx, lib.Int64ToString(userData.ID), parentID, dirName)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}

	dir := entity.Dir{}
	dir.DirName = dirName
	dir.ParentID = parentIDInt
	dir.UserID = userData.ID

	if idInt > 0 {
		b, s, dirData := model.DirData(tx, id)
		if !b {
			tx.Rollback()
			res.Message = s
			return res
		}
		if dirData.ID == 0 {
			tx.Rollback()
			res.Message = lang.DirectoryDoesNotExist
			return res
		}
		if dirData.UserID != userData.ID {
			tx.Rollback()
			res.Message = lang.NoPermission
			return res
		}
		if dirData.ID == parentIDInt {
			tx.Rollback()
			res.Message = lang.Typo
			return res
		}

		if sd.ID > 0 && sd.ID != idInt && sd.DirName == dirName {
			tx.Rollback()
			res.Message = lang.DirectoryAlreadyExists
			return res
		}

		b, s, r := model.DirUpdate(tx, id, dir)
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

		lib.WriteLog(userData.Account, "modify folder id: "+id)

		res.State = true
		res.Data = r
	} else {
		if sd.ID > 0 {
			tx.Rollback()
			res.Message = lang.DirectoryAlreadyExists
			return res
		}

		b, s, r := model.DirAdd(tx, dir)
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

		lib.WriteLog(userData.Account, "new folder "+dirName)

		res.State = true
		res.Data = r
	}

	tx.Commit()
	db.Close()
	return res
}

func DirInfo(userToken, id string) entity.Result {
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

	_, _, tx, db := model.ConnDB()

	b, s, r := model.DirData(tx, id)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}
	if r.ID == 0 {
		tx.Rollback()
		res.Message = lang.DirectoryDoesNotExist
		return res
	}
	if r.ID > 0 && r.UserID != userData.ID {
		tx.Rollback()
		res.Message = lang.NoPermission
		return res
	}

	res.State = true
	res.Data = r

	tx.Commit()
	db.Close()
	return res
}

func DirDel(userToken, id string) entity.Result {
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

	if id == "" {
		res.Message = lang.Typo
		return res
	}

	_, _, tx, db := model.ConnDB()

	b, s, r := model.DirData(tx, id)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}
	if r.ID == 0 {
		tx.Rollback()
		res.Message = lang.DirectoryDoesNotExist
		return res
	}
	if r.UserID != userData.ID {
		tx.Rollback()
		res.Message = lang.NoPermission
		return res
	}

	b, s = dirRec(tx, id, userData.ID)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}

	lib.WriteLog(userData.Account, "delete folder "+r.DirName)

	res.State = true
	tx.Commit()
	db.Close()

	return res
}

func dirRec(tx *sql.Tx, id string, userID int64) (bool, string) {
	lang := lang.Lang()

	_, _, idInt := lib.StringToInt64(id)

	fileList := model.Files(tx, 0, "", userID, idInt, 0)
	if len(fileList) > 0 {
		for i := 0; i < len(fileList); i++ {
			if fileList[i].Status == 1 {
				return false, lang.AbnormalFileStatus
			}
		}
	}

	// 删除文件
	b, s, _ := model.FileDelDir(tx, id)
	if !b {
		return false, s
	}

	// 删除文件夹
	dirList := model.Dirs(tx, 0, "", idInt, userID)
	if len(dirList) > 0 {
		for i := 0; i < len(dirList); i++ {
			b, s := dirRec(tx, lib.IntToString(int(dirList[i].ID)), userID)
			if !b {
				return b, s
			}
		}
	}

	b, s, _ = model.DirDel(tx, id)
	if !b {
		return b, s
	}

	return true, ""
}

func DirMove(userToken, id, ids string) entity.Result {
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

	if id == "" {
		res.Message = lang.Typo
		return res
	}

	if ids == "" {
		res.Message = lang.Typo
		return res
	}

	_, _, tx, db := model.ConnDB()

	_, _, idInt := lib.StringToInt(id)
	if idInt > 0 {
		b, s, r := model.DirData(tx, id)
		if !b {
			tx.Rollback()
			res.Message = s
			return res
		}
		if r.ID == 0 {
			tx.Rollback()
			res.Message = lang.DirectoryDoesNotExist
			return res
		}
		if r.UserID != userData.ID {
			tx.Rollback()
			res.Message = lang.NoPermission
			return res
		}
	}

	logInfo := ""
	if lib.StringContains(ids, ",") {
		idArr := strings.Split(ids, ",")
		for i := 0; i < len(idArr); i++ {
			if idArr[i] == id {
				tx.Rollback()
				res.Message = lang.Typo
				return res
			}
		}
		for i := 0; i < len(idArr); i++ {
			b, s, r := model.DirData(tx, idArr[i])
			if !b {
				tx.Rollback()
				res.Message = s
				return res
			}
			if r.ID == 0 {
				tx.Rollback()
				res.Message = lang.DirectoryDoesNotExist
				return res
			}
			if r.UserID != userData.ID {
				tx.Rollback()
				res.Message = lang.NoPermission
				return res
			}
			logInfo += "," + r.DirName
		}
	} else {
		if ids == id {
			tx.Rollback()
			res.Message = lang.Typo
			return res
		}
		b, s, r := model.DirData(tx, ids)
		if !b {
			tx.Rollback()
			res.Message = s
			return res
		}
		if r.ID == 0 {
			tx.Rollback()
			res.Message = lang.DirectoryDoesNotExist
			return res
		}
		if r.UserID != userData.ID {
			tx.Rollback()
			res.Message = lang.NoPermission
			return res
		}
		logInfo = r.DirName
	}

	b, s, _ := model.DirMove(tx, id, ids)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}

	lib.WriteLog(userData.Account, "move folder "+logInfo)

	res.State = true
	tx.Commit()
	db.Close()

	return res
}
