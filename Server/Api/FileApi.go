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

func FileModify(w http.ResponseWriter, r *http.Request) {
	userToken := strings.TrimSpace(Post(r, "userToken"))
	id := strings.TrimSpace(Post(r, "id"))
	fileName := strings.TrimSpace(Post(r, "fileName"))
	dirID := strings.TrimSpace(Post(r, "dirID"))
	HttpWrite(w, processmodule.FileModify(userToken, id, fileName, dirID))
}

func Files(w http.ResponseWriter, r *http.Request) {
	userToken := strings.TrimSpace(Post(r, "userToken"))
	order := strings.TrimSpace(Post(r, "order"))
	fileName := strings.TrimSpace(Post(r, "fileName"))
	dirID := strings.TrimSpace(Post(r, "dirID"))
	HttpWrite(w, processmodule.Files(userToken, order, fileName, dirID))
}

func FileDel(w http.ResponseWriter, r *http.Request) {
	userToken := strings.TrimSpace(Post(r, "userToken"))
	id := strings.TrimSpace(Post(r, "id"))
	HttpWrite(w, processmodule.FileDel(userToken, id))
}

func FileMove(w http.ResponseWriter, r *http.Request) {
	userToken := strings.TrimSpace(Post(r, "userToken"))
	dirID := strings.TrimSpace(Post(r, "dirID"))
	ids := strings.TrimSpace(Post(r, "ids"))
	HttpWrite(w, processmodule.FileMove(userToken, dirID, ids))
}

func FileUpload(w http.ResponseWriter, r *http.Request) {
	// userToken := strings.TrimSpace(Post(r, "userToken"))
	// id := strings.TrimSpace(Post(r, "id"))
	// HttpWrite(w, processmodule.FileUpload(userToken, id))
}

func FileDownload(w http.ResponseWriter, r *http.Request) {
	// userToken := strings.TrimSpace(Post(r, "userToken"))
	// id := strings.TrimSpace(Post(r, "id"))
	// HttpWrite(w, processmodule.FileDownload(userToken, id))
}
