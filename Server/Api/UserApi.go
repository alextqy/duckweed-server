package api

import (
	processmodule "duckweed-server/Server/ProcessModule"
	"net/http"
)

func UserList(w http.ResponseWriter, r *http.Request) {
	page := Post(r, "page")
	pageSize := Post(r, "pageSize")
	order := Post(r, "order")
	account := Post(r, "account")
	name := Post(r, "name")
	level := Post(r, "level")
	status := Post(r, "status")
	HttpWrite(w, processmodule.UserList(page, pageSize, order, account, name, level, status))
}

func UserGet(w http.ResponseWriter, r *http.Request) {
	HttpWrite(w, processmodule.UserGet(Get(r, "id")))
}

func UserCheck(w http.ResponseWriter, r *http.Request) {
	account := Post(r, "account")
	name := Post(r, "name")
	password := Post(r, "password")
	level := Post(r, "level")
	id := Post(r, "id")
	HttpWrite(w, processmodule.UserCheck(account, name, password, level, id))
}

func UserDel(w http.ResponseWriter, r *http.Request) {
	HttpWrite(w, processmodule.UserDel(Get(r, "id")))
}
