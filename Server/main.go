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
	if !lib.FileExist("./Dao.db") {
		_, memo := lib.FileMake("./Dao.db")
		if memo != "" {
			log.Fatal(memo)
		}
	}

	port := ""
	fmt.Println("Enter port number: ")
	fmt.Scan(&port)

	mux := http.NewServeMux()
	routes(mux)
	server := &http.Server{
		Addr:         ":" + port,
		WriteTimeout: time.Second * 5, //设置写超时
		ReadTimeout:  time.Second * 5, //设置读超时
		Handler:      mux,
	}
	log.Println("http server on port:" + port)
	log.Fatal(server.ListenAndServe())
}

func routes(mux *http.ServeMux) {
	mux.HandleFunc("/test", api.Test)
}
