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

func SetAvailableSpace(w http.ResponseWriter, r *http.Request) {
	userToken := strings.TrimSpace(Post(r, "userToken"))
	id := strings.TrimSpace(Post(r, "id"))
	availableSpace := strings.TrimSpace(Post(r, "availableSpace"))
	HttpWrite(w, processmodule.SetAvailableSpace(userToken, id, availableSpace))
}

func SetRootAccount(w http.ResponseWriter, r *http.Request) {
	userToken := strings.TrimSpace(Post(r, "userToken"))
	id := strings.TrimSpace(Post(r, "id"))
	HttpWrite(w, processmodule.SetRootAccount(userToken, id))
}

func DisableUser(w http.ResponseWriter, r *http.Request) {
	userToken := strings.TrimSpace(Post(r, "userToken"))
	id := strings.TrimSpace(Post(r, "id"))
	HttpWrite(w, processmodule.DisableUser(userToken, id))
}

func UserDel(w http.ResponseWriter, r *http.Request) {
	userToken := strings.TrimSpace(Post(r, "userToken"))
	id := strings.TrimSpace(Post(r, "id"))
	HttpWrite(w, processmodule.UserDel(userToken, id))
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	account := strings.TrimSpace(Post(r, "account"))
	name := strings.TrimSpace(Post(r, "name"))
	password := strings.TrimSpace(Post(r, "password"))
	email := strings.TrimSpace(Post(r, "email"))
	captcha := strings.TrimSpace(Post(r, "captcha"))
	HttpWrite(w, processmodule.SignUp(account, name, password, email, captcha))
}

func CheckPersonalData(w http.ResponseWriter, r *http.Request) {
	userToken := strings.TrimSpace(Post(r, "userToken"))
	HttpWrite(w, processmodule.CheckPersonalData(userToken))
}

func ModifyPersonalData(w http.ResponseWriter, r *http.Request) {
	userToken := strings.TrimSpace(Post(r, "userToken"))
	name := strings.TrimSpace(Post(r, "name"))
	password := strings.TrimSpace(Post(r, "password"))
	email := strings.TrimSpace(Post(r, "email"))
	captcha := strings.TrimSpace(Post(r, "captcha"))
	HttpWrite(w, processmodule.ModifyPersonalData(userToken, name, password, email, captcha))
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	newPassword := strings.TrimSpace(Post(r, "newPassword"))
	captcha := strings.TrimSpace(Post(r, "captcha"))
	HttpWrite(w, processmodule.ResetPassword(newPassword, captcha))
}

func SendEmail(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(Post(r, "email"))
	sendType := strings.TrimSpace(Post(r, "sendType"))
	HttpWrite(w, processmodule.SendEmail(email, sendType))
}
