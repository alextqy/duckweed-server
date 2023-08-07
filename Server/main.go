package main

import (
	api "duckweed-server/Server/Api"
	lib "duckweed-server/Server/Lib"
	model "duckweed-server/Server/Model"
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
	var confModel model.ConfModel
	_, byteData := lib.FileRead("./Conf.json")
	json.Unmarshal([]byte(byteData), &confModel)

	// 新建数据库文件
	if !lib.FileExist("./Dao.db") {
		_, memo := lib.FileMake("./Dao.db")
		if memo != "" {
			panic(memo)
		}
	}

	// 开启内网广播
	go loopBroadcast(ips[len(ips)-1], confModel.UdpPort)

	mux := http.NewServeMux()
	routes(mux)
	server := &http.Server{
		Addr:         ":" + confModel.TcpPort,
		WriteTimeout: time.Second * 5, //设置写超时
		ReadTimeout:  time.Second * 5, //设置读超时
		Handler:      mux,
	}
	log.Println("Http server on port:" + confModel.TcpPort)
	log.Fatal(server.ListenAndServe())
}

func routes(mux *http.ServeMux) {
	mux.HandleFunc("/test", api.Test)
}

func loopBroadcast(ip string, port string) {
	for {
		lib.Broadcast(port, ip+":"+port)
		time.Sleep(time.Second)
	}
}
