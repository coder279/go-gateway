package main

import (
	"bufio"
	"log"
	"net/http"
	"net/url"
)
var (
	proxy_addr = "http://127.0.0.1:2003"
	port       = "2002"
)

func handler(w http.ResponseWriter,req *http.Request){
	//step 1 解析代理地址 更改请求体内的协议和主机
	proxy,err := url.Parse(proxy_addr)
	req.URL.Scheme = proxy.Scheme
	req.URL.Host = proxy.Host
	//step 2 请求下游
	transport := http.DefaultTransport
	resp,err := transport.RoundTrip(req)
	if err != nil {
		log.Println(err)
		return
	}
	//step 3 把下游请求内容返回给上游
	for k,vv := range resp.Header {
		for _,v:= range vv {
			w.Header().Add(k,v)
		}
	}
	defer resp.Body.Close()
	bufio.NewReader(resp.Body).WriteTo(w)

}

func main() {
	http.HandleFunc("/",handler)
	log.Println("start server now...")
	err := http.ListenAndServe(":"+port,nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}
