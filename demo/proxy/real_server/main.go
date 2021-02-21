package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	rs1 := &RealSever{Addr:"127.0.0.1:2003"}
	rs1.Run()
	rs2 := &RealSever{Addr:"127.0.0.1:2004"}
	rs2.Run()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //监听两种信号形式
	<-quit
}
type RealSever struct {
	Addr string
}
//启动服务
func (r *RealSever)Run(){
	log.Println("Starting httpserver at" + r.Addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/",r.HelloHandler)
	mux.HandleFunc("/base/error", r.ErrorHandler)
	mux.HandleFunc("/test_http_string/test_http_string/aaa", r.TimeoutHandler)
	server := &http.Server{
		Addr:r.Addr,
		WriteTimeout: 10 * time.Second,
		Handler: mux,
	}
	go func() {
		log.Fatal(server.ListenAndServe())
	}()
}
func (r *RealSever) HelloHandler(w http.ResponseWriter,req *http.Request){
	uppath := fmt.Sprintf("http://%s%s\n",r.Addr,req.URL.Path) //路由
	//remote_addr 客户端ip
	//x-forwarded-for 如果用了反向代理 用户会将代理的ip写入这个里面

	realIp := fmt.Sprintf("RemoteAddr=%x,X-Forwarded-For=%v,X-RealIp=%v\n",req.RemoteAddr,req.Header.Get("X-Forwarded-For"),req.Header.Get("X-Real-Ip"))
	header := fmt.Sprintf("header=%v\n",req.Header)
	io.WriteString(w, uppath)
	io.WriteString(w, realIp)
	io.WriteString(w, header)
}
func (r *RealSever) ErrorHandler(w http.ResponseWriter, req *http.Request) {
	upath := "error handler"
	w.WriteHeader(500)
	io.WriteString(w, upath)
}

func (r *RealSever) TimeoutHandler(w http.ResponseWriter, req *http.Request) {
	time.Sleep(6*time.Second)
	upath := "timeout handler"
	w.WriteHeader(200)
	io.WriteString(w, upath)
}
