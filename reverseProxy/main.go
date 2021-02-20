package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var addr = "127.0.0.1:2004"
func main() {
	rs := "http://www.baidu.com"
	url,err := url.Parse(rs)
	if err != nil {
		log.Fatalf("转发失败")
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	log.Println("starting httpserve at"+ addr)
	log.Fatal(http.ListenAndServe(addr,proxy))

}
func HandlerHttp(w *http.ResponseWriter, r *http.Request){
	fmt.Println("开心呢")
}

