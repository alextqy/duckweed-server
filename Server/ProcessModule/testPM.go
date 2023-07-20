package processmodule

import entity "duckweed-server/Server/Entity"

func Test() *entity.Result {
	res := new(entity.Result)
	res.State = true
	res.Code = 200
	res.Message = "succeed"
	res.Data = ""
	return res
}
