package processmodule

import (
	entity "duckweed-server/Server/Entity"
)

func Test(text string) entity.Result {
	// _, _, encrypted := lib.AesEncrypterCBC(text, "aaaaaaaaaaaaaaaa", "bbbbbbbbbbbbbbbb")
	// _, _, decrypted := lib.AesDecrypterCBC(encrypted, "aaaaaaaaaaaaaaaa", "bbbbbbbbbbbbbbbb")
	// compile := lib.RegEnNum(text)
	res := entity.Result{
		State:   true,
		Code:    200,
		Message: "succeed",
		Data:    nil,
	}
	return res
}
