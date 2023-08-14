package api

import (
	processmodule "duckweed-server/Server/ProcessModule"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request) {
	HttpWrite(w, processmodule.Test(Post(r, "text")))
}
