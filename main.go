package main

import (
	//"fmt"
	"log"
	"net/http"

	//"wxcloudrun-golang/db"
	"wxcloudrun-golang/service"
)

func main() {
	//if err := db.Init(); err != nil {
	//	panic(fmt.Sprintf("mysql init failed with %+v", err))
	//}

	http.HandleFunc("/", service.IndexHandler)
	http.HandleFunc("/api/count", service.CounterHandler)
	http.HandleFunc("/api/hello", service.HelloHandler)
	http.HandleFunc("/api/message", service.WXMessageHandler)

	log.Fatal(http.ListenAndServe(":3000", nil))
}
