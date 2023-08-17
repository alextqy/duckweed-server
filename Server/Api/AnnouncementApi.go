package api

import (
	processmodule "duckweed-server/Server/ProcessModule"
	"net/http"
	"strings"
)

func Announcements(w http.ResponseWriter, r *http.Request) {
	HttpWrite(w, processmodule.Announcements(strings.TrimSpace(Post(r, "userToken"))))
}

func AnnouncementGet(w http.ResponseWriter, r *http.Request) {
	userToken := strings.TrimSpace(Post(r, "userToken"))
	id := strings.TrimSpace(Post(r, "id"))
	HttpWrite(w, processmodule.AnnouncementGet(userToken, id))
}

func AnnouncementAdd(w http.ResponseWriter, r *http.Request) {
	userToken := strings.TrimSpace(Post(r, "userToken"))
	content := strings.TrimSpace(Post(r, "content"))
	HttpWrite(w, processmodule.AnnouncementAdd(userToken, content))
}

func AnnouncementDel(w http.ResponseWriter, r *http.Request) {
	userToken := strings.TrimSpace(Post(r, "userToken"))
	id := strings.TrimSpace(Post(r, "id"))
	HttpWrite(w, processmodule.AnnouncementDel(userToken, id))
}
