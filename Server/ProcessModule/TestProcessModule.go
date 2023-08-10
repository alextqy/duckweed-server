package processmodule

import entity "duckweed-server/Server/Entity"

func Test() entity.Result {
	res := entity.Result{
		State:   true,
		Code:    200,
		Message: "succeed",
		Data:    "",
	}
	return res
}
