package main

import (
	"log"
	"net/http"
	"time"
)

var Addr = ":1210"
func main() {
	//创建路由
	mux := http.NewServeMux()
	//设置路由匹配方法
	mux.HandleFunc("/bye",sayBye)
	//创建服务器
	server := &http.Server{
		Addr: Addr,
		WriteTimeout:10 * time.Second, //请求超时时间
		Handler:mux,
	}
	log.Println("starting at server" + Addr)
	log.Println(server.ListenAndServe())
}
func sayBye(w http.ResponseWriter,r *http.Request){
	time.Sleep(5 * time.Second)
	w.Write([]byte("byte bye"))
}
