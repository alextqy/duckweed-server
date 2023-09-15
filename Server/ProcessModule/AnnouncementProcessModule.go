package processmodule

import (
	entity "duckweed-server/Server/Entity"
	lang "duckweed-server/Server/Lang"
	lib "duckweed-server/Server/Lib"
	model "duckweed-server/Server/Model"
)

func Announcements(userToken string) entity.Result {
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

	res.State = true
	res.Code = 200
	res.Message = ""
	res.Data = model.Announcements(tx, -1)

	tx.Commit()
	db.Close()
	return res
}

func AnnouncementGet(userToken, id string) entity.Result {
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
	b, s, r := model.AnnouncementData(tx, id)
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

func AnnouncementAdd(userToken, content string) entity.Result {
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

	if content == "" {
		res.Message = lang.Typo
		return res
	}

	if !lib.RegWriting(content) {
		res.Message = lang.MalformedContent
		return res
	}

	_, _, tx, db := model.ConnDB()
	data := entity.AnnouncementEntity{}
	data.Content = content
	b, s, r := model.AnnouncementAdd(tx, data)
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

	lib.WriteLog(adminAccount, adminAccount+" new announcement content: "+content)

	tx.Commit()
	res.Message = ""
	res.State = true
	res.Data = r

	db.Close()
	return res
}

func AnnouncementDel(userToken, id string) entity.Result {
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

	_, _, tx, db := model.ConnDB()

	b, s, data := model.AnnouncementData(tx, id)
	if !b {
		tx.Rollback()
		res.Message = s
		return res
	}
	if data.ID == 0 {
		tx.Rollback()
		res.Message = lang.NoData
		return res
	}

	b, s, r := model.AnnouncementDel(tx, id)
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

	lib.WriteLog(adminAccount, adminAccount+" delete announcement data content: "+data.Content)

	tx.Commit()
	res.Message = ""
	res.State = true
	res.Data = r

	db.Close()
	return res
}
