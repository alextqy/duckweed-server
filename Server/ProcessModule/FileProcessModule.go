package processmodule

import (
	entity "duckweed-server/Server/Entity"
	lang "duckweed-server/Server/Lang"
	lib "duckweed-server/Server/Lib"
	model "duckweed-server/Server/Model"
	"strings"
)

func FileAdd(userToken, fileName, fileType, fileSize, md5, dirID, sourceAddress string) entity.Result {
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
	if sourceAddress == "" {
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
	file.SourceAddress = sourceAddress

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

func FileModify(userToken, id, fileName, dirID, sourceAddress string) entity.Result {
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
	if fileData.UserID != userData.ID {
		tx.Rollback()
		res.Message = lang.NoPermission
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
	fileData.SourceAddress = sourceAddress
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

func Files(userToken, order, fileName, dirID, status string) entity.Result {
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
	_, _, dirIDInt := lib.StringToInt(dirID)
	_, _, statusInt := lib.StringToInt(status)
	if dirIDInt < 0 {
		res.Message = lang.Typo
		return res
	}
	if statusInt < 0 {
		res.Message = lang.Typo
		return res
	}

	_, _, tx, db := model.ConnDB()

	res.State = true
	res.Data = model.Files(tx, orderInt, fileName, userData.ID, dirIDInt, statusInt)

	tx.Commit()
	db.Close()
	return res
}

func FileDel(userToken, id string) entity.Result {
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

	logContent := ""
	if lib.StringContains(id, ",") {
		idArr := strings.Split(id, ",")
		for i := 0; i < len(idArr); i++ {
			b, s, fileData := model.FileData(tx, idArr[i])
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
			if fileData.UserID != userData.ID {
				tx.Rollback()
				res.Message = lang.NoPermission
				return res
			}
		}

		logContent = " delete file data id: " + id
	} else {
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
		if fileData.UserID != userData.ID {
			tx.Rollback()
			res.Message = lang.NoPermission
			return res
		}

		if fileData.FileType != "" {
			fileData.FileType = "." + fileData.FileType
		}

		logContent = " delete file data " + fileData.FileName + fileData.FileType
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

	lib.WriteLog(userData.Account, logContent)

	res.State = true

	tx.Commit()
	db.Close()
	return res
}

func FileMove(userToken, dirID, ids string) entity.Result {
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

	if dirID == "" {
		res.Message = lang.Typo
		return res
	}

	if ids == "" {
		res.Message = lang.Typo
		return res
	}

	_, _, tx, db := model.ConnDB()

	b, s, r := model.DirData(tx, dirID)
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

	logInfo := ""
	if lib.StringContains(ids, ",") {
		idArr := strings.Split(ids, ",")
		for i := 0; i < len(idArr); i++ {
			b, s, r := model.FileData(tx, idArr[i])
			if !b {
				tx.Rollback()
				res.Message = s
				return res
			}
			if r.ID == 0 {
				tx.Rollback()
				res.Message = lang.FileDoesNotExist
				return res
			}
			if r.UserID != userData.ID {
				tx.Rollback()
				res.Message = lang.NoPermission
				return res
			}
			logInfo += "," + r.FileName
		}
	} else {
		b, s, r := model.FileData(tx, ids)
		if !b {
			tx.Rollback()
			res.Message = s
			return res
		}
		if r.ID == 0 {
			tx.Rollback()
			res.Message = lang.FileDoesNotExist
			return res
		}
		if r.UserID != userData.ID {
			tx.Rollback()
			res.Message = lang.NoPermission
			return res
		}
		logInfo = r.FileName
	}

	b, s, _ = model.FileMove(tx, dirID, ids)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}

	lib.WriteLog(userData.Account, "move file "+logInfo)

	res.State = true
	tx.Commit()
	db.Close()
	return res
}

// _, f := lib.Filespec("../Temp/Upload/sqlite4.dll")
// _, _, r := lib.FileReadBlock("../Temp/Upload/sqlite3.dll", 10, 11)
// b, s := lib.FileWriteByte("../Temp/Upload/sqlite4.dll", r)
