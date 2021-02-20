package main

import (
	"bufio"
	"log"
	"net/http"
	"net/url"
)

var (
	proxy_addr = "http://127.0.0.1:2003" //代理的地址
	port = "2002"
)
func handler(w http.ResponseWriter,r *http.Request){
	//step1 解析代理地址并修改更改请求体的协议和主机
	proxy,err := url.Parse(proxy_addr)
	r.URL.Scheme = proxy.Scheme //协议
	r.URL.Host = proxy.Host     //主机
	//step2 请求下游
	transport := http.DefaultTransport
	resp,err := transport.RoundTrip(r)
	if err != nil {
		log.Print(err)
		return
	}
	//step3 把下游请求内容返回给上游
	for k,vv := range resp.Header{
		for _,v := range vv {
			w.Header().Add(k,v)
		}
	}
	defer resp.Body.Close()
	//将下游信息写给上游
	bufio.NewReader(resp.Body).WriteTo(w)
}
func main() {
	http.HandleFunc("/",handler)
	log.Println("start serving on port" +port)
	err := http.ListenAndServe(":"+port,nil)
	if err != nil {
		log.Fatal(err)
	}
}
