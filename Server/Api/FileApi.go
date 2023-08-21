package api

import (
	processmodule "duckweed-server/Server/ProcessModule"
	"net/http"
	"strings"
)

func FileAdd(w http.ResponseWriter, r *http.Request) {
	userToken := strings.TrimSpace(Post(r, "userToken"))
	fileName := strings.TrimSpace(Post(r, "fileName"))
	fileType := strings.TrimSpace(Post(r, "fileType"))
	fileSize := strings.TrimSpace(Post(r, "fileSize"))
	md5 := strings.TrimSpace(Post(r, "md5"))
	dirID := strings.TrimSpace(Post(r, "dirID"))
	HttpWrite(w, processmodule.FileAdd(userToken, fileName, fileType, fileSize, md5, dirID))
}
