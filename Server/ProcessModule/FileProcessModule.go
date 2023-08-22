package processmodule

import (
	entity "duckweed-server/Server/Entity"
	lang "duckweed-server/Server/Lang"
	lib "duckweed-server/Server/Lib"
	model "duckweed-server/Server/Model"
)

func FileAdd(userToken string, fileName string, fileType string, fileSize string, md5 string, dirID string) entity.Result {
	lang := lang.Lang()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
		Data:    nil,
	}

	if fileName == "" {
		res.Message = lang.Typo
		return res
	}
	if !lib.RegAll(fileName) {
		res.Message = lang.FileNameFormatError
		return res
	}
	if fileType != "" {
		if len(fileType) > 16 {
			res.Message = lang.Typo
			return res
		}
	}
	if fileSize == "" {
		res.Message = lang.Typo
		return res
	}
	if md5 == "" {
		res.Message = lang.Typo
		return res
	}
	if dirID == "" {
		res.Message = lang.Typo
		return res
	}

	_, _, fileSizeInt := lib.StringToInt(fileSize)
	_, _, dirIDInt := lib.StringToInt(dirID)
	if fileSizeInt <= 0 {
		res.Message = lang.Typo
		return res
	}
	if dirIDInt < 0 {
		res.Message = lang.Typo
		return res
	}

	userData := CheckToken(userToken)
	if userData.ID == 0 {
		res.Message = lang.ReLoginRequired
		return res
	}

	userSpace := "../Space/" + userData.Account + "/"
	outreachID := lib.MD5(lib.Int64ToString(lib.TimeStampMS()) + userData.Account + lib.RandStr(5))
	fileSpace := userSpace + outreachID

	_, _, tx, db := model.ConnDB()

	b, s, checkFile := model.FileDataSame(tx, dirID, fileName, fileType)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}
	if checkFile.ID > 0 {
		tx.Rollback()
		res.Message = lang.FileAlreadyExists
		return res
	}

	file := entity.FileEntity{}
	file.FileName = fileName
	file.FileType = fileType
	file.FileSize = fileSize
	file.StoragePath = fileSpace
	file.MD5 = md5
	file.UserID = userData.ID
	file.DirID = dirIDInt
	file.OutreachID = outreachID

	b, s, r := model.FileAdd(tx, file)
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

	b, s = lib.DirMake(fileSpace)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}

	if fileType != "" {
		fileType = "." + fileType
	}

	lib.WriteLog(userData.Account, "create a new file "+fileName+fileType)

	res.State = true
	res.Data = r

	tx.Commit()
	db.Close()
	return res
}

func FileModify(userToken string, id string, fileName string, dirID string) entity.Result {
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
	if fileName == "" {
		res.Message = lang.Typo
		return res
	}
	if !lib.RegAll(fileName) {
		res.Message = lang.FileNameFormatError
		return res
	}
	if dirID == "" {
		res.Message = lang.Typo
		return res
	}
	_, _, dirIDInt := lib.StringToInt(dirID)
	if dirIDInt < 0 {
		res.Message = lang.Typo
		return res
	}

	userData := CheckToken(userToken)
	if userData.ID == 0 {
		res.Message = lang.ReLoginRequired
		return res
	}

	_, _, tx, db := model.ConnDB()

	b, s, fileData := model.FileData(tx, id)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}
	if fileData.ID == 0 {
		tx.Rollback()
		res.Message = lang.NoData
		return res
	}
	if dirID != "" && dirIDInt > 0 {
		b, s, dirData := model.DirData(tx, dirID)
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
		fileData.DirID = dirIDInt
	}

	fileData.FileName = fileName
	b, s, r := model.FileUpdate(tx, id, fileData)
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

	lib.WriteLog(userData.Account, "modify file id: "+id)

	res.State = true
	res.Data = r

	tx.Commit()
	db.Close()
	return res
}

func Files(userToken string, order string, fileName string, dirID string) entity.Result {
	lang := lang.Lang()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
		Data:    nil,
	}

	_, _, orderInt := lib.StringToInt(order)
	_, _, dirIDInt := lib.StringToInt(dirID)
	if dirIDInt < 0 {
		res.Message = lang.Typo
		return res
	}

	userData := CheckToken(userToken)
	if userData.ID == 0 {
		res.Message = lang.ReLoginRequired
		return res
	}

	_, _, tx, db := model.ConnDB()

	res.State = true
	res.Data = model.Files(tx, orderInt, fileName, userData.ID, dirIDInt)

	tx.Commit()
	db.Close()
	return res
}

func FileDel(userToken string, id string) entity.Result {
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

	userData := CheckToken(userToken)
	if userData.ID == 0 {
		res.Message = lang.ReLoginRequired
		return res
	}

	_, _, tx, db := model.ConnDB()

	b, s, fileData := model.FileData(tx, id)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}
	if fileData.ID == 0 {
		tx.Rollback()
		res.Message = lang.NoData
		return res
	}

	if fileData.FileType != "" {
		fileData.FileType = "." + fileData.FileType
	}

	b, s, r := model.FileDel(tx, id)
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

	lib.WriteLog(userData.Account, " delete file data "+fileData.FileName+fileData.FileType)

	res.State = true

	tx.Commit()
	db.Close()
	return res
}

// _, f := lib.Filespec("../Temp/Upload/sqlite4.dll")
// _, _, r := lib.FileReadBlock("../Temp/Upload/sqlite3.dll", 10, 11)
// b, s := lib.FileWriteByte("../Temp/Upload/sqlite4.dll", r)
