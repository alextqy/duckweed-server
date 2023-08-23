package processmodule

import (
	entity "duckweed-server/Server/Entity"
	lang "duckweed-server/Server/Lang"
	lib "duckweed-server/Server/Lib"
)

func ViewLog(userToken, date, account string) entity.Result {
	lang := lang.Lang()
	res := entity.Result{
		State:   false,
		Code:    200,
		Message: "",
		Data:    nil,
	}

	if date == "" {
		res.Message = lang.Typo
		return res
	}
	if account == "" {
		res.Message = lang.Typo
		return res
	}

	permissions, _ := CheckLevel(userToken)
	if permissions != 2 {
		res.Message = lang.NoPermission
		return res
	}

	logFile := "../Log/" + date + "/" + account + ".log"

	b := lib.FileExist(logFile)
	if !b {
		res.Message = lang.LogDoesNotExist
		return res
	}

	b, s := lib.FileRead(logFile)
	if !b {
		res.Message = s
		return res
	}

	res.State = true
	res.Data = s

	return res
}
