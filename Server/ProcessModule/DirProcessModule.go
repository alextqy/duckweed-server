package processmodule

import (
	entity "duckweed-server/Server/Entity"
	lang "duckweed-server/Server/Lang"
	lib "duckweed-server/Server/Lib"
	model "duckweed-server/Server/Model"
)

func Dirs(userToken string, order string, parentID string, dirName string) entity.Result {
	lang := lang.Lang()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
		Data:    nil,
	}

	if parentID == "" {
		res.Message = lang.Typo
		return res
	}

	_, _, orderInt := lib.StringToInt(order)
	_, _, parentIDInt := lib.StringToInt(parentID)
	if parentIDInt < 0 {
		res.Message = lang.Typo
		return res
	}

	userData := CheckToken(userToken)
	if userData.ID == 0 {
		res.Message = lang.NoData
		return res
	}

	_, _, tx, db := model.ConnDB()

	res.State = true
	res.Data = model.Dirs(tx, orderInt, dirName, parentIDInt, userData.ID)

	tx.Commit()
	db.Close()
	return res
}

func DirAdd(userToken string, dirName string, parentID string) entity.Result {
	lang := lang.Lang()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
		Data:    nil,
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
	_, _, parentIDInt := lib.StringToInt(parentID)
	if parentIDInt < 0 {
		res.Message = lang.Typo
		return res
	}

	userData := CheckToken(userToken)
	if userData.ID == 0 {
		res.Message = lang.NoData
		return res
	}

	_, _, tx, db := model.ConnDB()

	if parentIDInt > 0 {
		b, s, r := model.DirData(tx, parentID)
		if !b {
			res.Message = s
			return res
		}
		if r.ID == 0 {
			res.Message = lang.ParentFolderDoesNotExist
			return res
		}
	}

	b, s, sd := model.DirDataSame(tx, dirName, parentID)
	if !b {
		res.Message = s
		return res
	}
	if sd.ID > 0 {
		res.Message = lang.DirectoryAlreadyExists
		return res
	}

	dir := entity.DirEntity{}
	dir.DirName = dirName
	dir.ParentID = parentIDInt
	dir.UserID = userData.ID
	b, s, r := model.DirAdd(tx, dir)
	if !b {
		res.Message = s
		return res
	}
	if r == 0 {
		res.Message = lang.OperationFailed
		return res
	}

	lib.WriteLog(userData.Account, "new folder "+dirName)

	res.State = true
	res.Data = r

	tx.Commit()
	db.Close()
	return res
}
