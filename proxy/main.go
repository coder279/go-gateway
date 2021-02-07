package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type Pxy struct {}

func (p *Pxy) ServeHTTP(rw http.ResponseWriter,req *http.Request){
	fmt.Printf("Received request %v,%v,%v\n",req.Method,req.Host,req.RemoteAddr)
	/*
		DefaultTransport is the default implementation of Transport and is used by DefaultClient.
		It establishes network connections as needed and caches them for reuse by subsequent calls.
		It uses HTTP proxies as directed by the $HTTP_PROXY and $NO_PROXY (or $http_proxy and $no_proxy) environment variables.
	*/
	transport := http.DefaultTransport
	//step 1 浅拷贝对象，然后就新增数据
	outReq := new(http.Request)
	*outReq = *req
	if clientIp,_,err := net.SplitHostPort(req.RemoteAddr);err==nil{
		//X-Forwarded-For（XFF）是用来识别通过HTTP代理或负载均衡方式连接到Web服务器的客户端最原始的IP地址的HTTP请求头字段。
		if prior,ok := outReq.Header["X-Forwarded-For"];ok{
			clientIp = strings.Join(prior,",") + "," + clientIp
		}
		//设置HTTP头部
		outReq.Header.Set("X-Forwarded-For",clientIp)
		//step2 请求下游
		//RoundTrip 代表一个http事务，给一个请求返回一个响应)
		//说白了，就是你给它一个request,它给你一个response
		res ,err := transport.RoundTrip(outReq)
		if err != nil {
			//可以做Header头部处理
			rw.WriteHeader(http.StatusBadGateway)
			return
		}
		for key,value := range res.Header{
			for _,v := range value{
				rw.Header().Add(key,v)
			}
		}
		rw.WriteHeader(res.StatusCode)
		io.Copy(rw,res.Body)
		res.Body.Close()
	}
}
func main() {
	fmt.Println("Serve on :8080")
	http.Handle("/",&Pxy{})
	http.ListenAndServe("0.0.0.0:8080",nil)
}
