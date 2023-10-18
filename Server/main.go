package main

import (
	api "duckweed-server/Server/Api"
	lib "duckweed-server/Server/Lib"
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

	// 新建数据库文件
	if !lib.FileExist("../Dao.db") {
		DaoState, memo := lib.FileMake("../Dao.db")
		if !DaoState {
			log.Fatal(memo)
		}
	}

	go loopBroadcast(ips[len(ips)-1], lib.CheckConf().UdpPort)
	go systemLog()
	go space()

	mux := http.NewServeMux()
	routes(mux)
	server := &http.Server{
		Addr:         ":" + lib.CheckConf().TcpPort,
		WriteTimeout: time.Second * 5, //设置写超时
		ReadTimeout:  time.Second * 5, //设置读超时
		Handler:      mux,
	}
	log.Println("Http server on port:" + lib.CheckConf().TcpPort)
	log.Fatal(server.ListenAndServe())
}

// 开启内网广播
func loopBroadcast(ip string, port string) {
	for {
		lib.Broadcast(port, ip+":"+lib.CheckConf().TcpPort)
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
	mux.HandleFunc("/user/get", api.UserGet)
	mux.HandleFunc("/set/available/space", api.SetAvailableSpace)
	mux.HandleFunc("/set/root/account", api.SetRootAccount)
	mux.HandleFunc("/disable/user", api.DisableUser)
	mux.HandleFunc("/user/del", api.UserDel)
	mux.HandleFunc("/sign/up", api.SignUp)
	mux.HandleFunc("/check/personal/data", api.CheckPersonalData)
	mux.HandleFunc("/modify/personal/data", api.ModifyPersonalData)
	mux.HandleFunc("/reset/password", api.ResetPassword)
	mux.HandleFunc("/send/email", api.SendEmail)

	mux.HandleFunc("/announcements", api.Announcements)
	mux.HandleFunc("/announcement/get", api.AnnouncementGet)
	mux.HandleFunc("/announcement/add", api.AnnouncementAdd)
	mux.HandleFunc("/announcement/del", api.AnnouncementDel)

	mux.HandleFunc("/view/log", api.ViewLog)

	mux.HandleFunc("/dirs", api.Dirs)
	mux.HandleFunc("/dir/action", api.DirAction)
	mux.HandleFunc("/dir/info", api.DirInfo)
	mux.HandleFunc("/dir/del", api.DirDel)

	mux.HandleFunc("/file/add", api.FileAdd)
	mux.HandleFunc("/file/modify", api.FileModify)
	mux.HandleFunc("/files", api.Files)
	mux.HandleFunc("/file/del", api.FileDel)
	mux.HandleFunc("/file/upload", api.FileUpload)     // x
	mux.HandleFunc("/file/download", api.FileDownload) // x
}
