package lang

import (
	entity "duckweed-server/Server/Entity"
	lib "duckweed-server/Server/Lib"
	"encoding/json"
)

func Lang() map[string]string {
	info := make(map[string]string)
	var confEntity entity.ConfEntity
	_, byteData := lib.FileRead("./Conf.json")
	json.Unmarshal([]byte(byteData), &confEntity)
	if confEntity.Lang == "zh" {
		info["NoData"] = "无数据"
		info["IncorrectPassword"] = "密码错误"
	} else if confEntity.Lang == "en" {
		info["NoData"] = "no data"
		info["IncorrectPassword"] = "incorrect password"
	} else {
		info["NoData"] = ""
		info["IncorrectPassword"] = ""
	}
	return info
}
