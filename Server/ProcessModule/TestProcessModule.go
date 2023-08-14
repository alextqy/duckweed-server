package processmodule

import (
	entity "duckweed-server/Server/Entity"
	lib "duckweed-server/Server/Lib"
)

func Test(text string) entity.Result {
	_, _, encrypted := lib.AesEncrypterCBC(text, "aaaaaaaaaaaaaaaa", "bbbbbbbbbbbbbbbb")
	_, _, decrypted := lib.AesDecrypterCBC(encrypted, "aaaaaaaaaaaaaaaa", "bbbbbbbbbbbbbbbb")
	res := entity.Result{
		State:   true,
		Code:    200,
		Message: "succeed",
		Data:    decrypted,
	}
	return res
}
