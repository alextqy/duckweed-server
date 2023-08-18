package api

import (
	processmodule "duckweed-server/Server/ProcessModule"
	"net/http"
	"strings"
)

func Dirs(w http.ResponseWriter, r *http.Request) {
	userToken := strings.TrimSpace(Post(r, "userToken"))
	order := strings.TrimSpace(Post(r, "order"))
	parentID := strings.TrimSpace(Post(r, "parentID"))
	dirName := strings.TrimSpace(Post(r, "dirName"))
	HttpWrite(w, processmodule.Dirs(userToken, order, parentID, dirName))
}

func DirAdd(w http.ResponseWriter, r *http.Request) {
	userToken := strings.TrimSpace(Post(r, "userToken"))
	dirName := strings.TrimSpace(Post(r, "dirName"))
	parentID := strings.TrimSpace(Post(r, "parentID"))
	HttpWrite(w, processmodule.DirAdd(userToken, dirName, parentID))
}
