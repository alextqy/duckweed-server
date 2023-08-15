package main

import (
	api "duckweed-server/Server/Api"
	entity "duckweed-server/Server/Entity"
	lib "duckweed-server/Server/Lib"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Local IP address:")
	_, _, ips := lib.LocalIP()
	for i := 0; i < len(ips); i++ {
		fmt.Println(ips[i])
	}

	// 读取配置文件
	var confEntity entity.ConfEntity
	_, byteData := lib.FileRead("./Conf.json")
	json.Unmarshal([]byte(byteData), &confEntity)

	// 新建数据库文件
	if !lib.FileExist("../Dao.db") {
		DaoState, memo := lib.FileMake("../Dao.db")
		if !DaoState {
			log.Fatal(memo)
		}
	}

	go loopBroadcast(ips[len(ips)-1], confEntity.UdpPort)
	go systemLog()
	go space()

	mux := http.NewServeMux()
	routes(mux)
	server := &http.Server{
		Addr:         ":" + confEntity.TcpPort,
		WriteTimeout: time.Second * 5, //设置写超时
		ReadTimeout:  time.Second * 5, //设置读超时
		Handler:      mux,
	}
	log.Println("Http server on port:" + confEntity.TcpPort)
	log.Fatal(server.ListenAndServe())
}

// 开启内网广播
func loopBroadcast(ip string, port string) {
	for {
		lib.Broadcast(port, ip+":"+port)
		time.Sleep(time.Second)
	}
}

// 系统日志
func systemLog() {
	for {
		if !lib.FileExist(lib.LogDir()) {
			lib.DirMake(lib.LogDir())
		}
		time.Sleep(time.Second)
	}
}

// 用户空间
func space() {
	for {
		if !lib.FileExist("../Space") {
			lib.DirMake("../Space")
		}
		time.Sleep(time.Second)
	}
}

func routes(mux *http.ServeMux) {
	mux.HandleFunc("/test", api.Test)
	mux.HandleFunc("/sign/in", api.SignIn)
	mux.HandleFunc("/sign/out", api.SignOut)
	mux.HandleFunc("/user/list", api.UserList)
	mux.HandleFunc("/users", api.Users)
}
