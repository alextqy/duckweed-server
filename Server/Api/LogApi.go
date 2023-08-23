package api

import (
	processmodule "duckweed-server/Server/ProcessModule"
	"net/http"
	"strings"
)

func ViewLog(w http.ResponseWriter, r *http.Request) {
	userToken := strings.TrimSpace(Post(r, "userToken"))
	date := strings.TrimSpace(Post(r, "date"))
	account := strings.TrimSpace(Post(r, "account"))
	HttpWrite(w, processmodule.ViewLog(userToken, date, account))
}
