package api

import (
	processmodule "duckweed-server/Server/ProcessModule"
	"net/http"
	"strings"
)

func UserLogin(w http.ResponseWriter, r *http.Request) {
	account := strings.TrimSpace(Post(r, "account"))
	password := strings.TrimSpace(Post(r, "password"))
	HttpWrite(w, processmodule.UserLogin(account, password))
}
