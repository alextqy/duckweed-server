package main

import (
	api "duckweed-server/Server/Api"
	"log"
	"net/http"
	"time"
)

func main() {
	port := "8000"
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
