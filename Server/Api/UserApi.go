package api

import (
	processmodule "duckweed-server/Server/ProcessModule"
	"net/http"
	"strings"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	account := strings.TrimSpace(Post(r, "account"))
	password := strings.TrimSpace(Post(r, "password"))
	HttpWrite(w, processmodule.SignIn(account, password))
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	HttpWrite(w, processmodule.SignOut(strings.TrimSpace(Post(r, "userToken"))))
}

func UserList(w http.ResponseWriter, r *http.Request) {
	userToken := strings.TrimSpace(Post(r, "userToken"))
	page := strings.TrimSpace(Post(r, "page"))
	pageSize := strings.TrimSpace(Post(r, "pageSize"))
	order := strings.TrimSpace(Post(r, "order"))
	account := strings.TrimSpace(Post(r, "account"))
	name := strings.TrimSpace(Post(r, "name"))
	level := strings.TrimSpace(Post(r, "level"))
	status := strings.TrimSpace(Post(r, "status"))
	HttpWrite(w, processmodule.UserList(userToken, page, pageSize, order, account, name, level, status))
}

func Users(w http.ResponseWriter, r *http.Request) {
	userToken := strings.TrimSpace(Post(r, "userToken"))
	order := strings.TrimSpace(Post(r, "order"))
	account := strings.TrimSpace(Post(r, "account"))
	name := strings.TrimSpace(Post(r, "name"))
	level := strings.TrimSpace(Post(r, "level"))
	status := strings.TrimSpace(Post(r, "status"))
	HttpWrite(w, processmodule.Users(userToken, order, account, name, level, status))
}

func UserGet(w http.ResponseWriter, r *http.Request) {
	userToken := strings.TrimSpace(Post(r, "userToken"))
	id := strings.TrimSpace(Post(r, "id"))
	HttpWrite(w, processmodule.UserGet(userToken, id))
}

func UserAction(w http.ResponseWriter, r *http.Request) {
	userToken := strings.TrimSpace(Post(r, "userToken"))
	account := strings.TrimSpace(Post(r, "account"))
	name := strings.TrimSpace(Post(r, "name"))
	password := strings.TrimSpace(Post(r, "password"))
	level := strings.TrimSpace(Post(r, "level"))
	availableSpace := strings.TrimSpace(Post(r, "availableSpace"))
	id := strings.TrimSpace(Post(r, "id"))
	HttpWrite(w, processmodule.UserAction(userToken, account, name, password, level, availableSpace, id))
}

func UserDel(w http.ResponseWriter, r *http.Request) {
	userToken := strings.TrimSpace(Post(r, "userToken"))
	id := strings.TrimSpace(Post(r, "id"))
	HttpWrite(w, processmodule.UserDel(userToken, id))
}
