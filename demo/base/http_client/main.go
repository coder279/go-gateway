package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func main() {
	//创建连接池
	transport := &http.Transport{
		//创建上下文呢
		DialContext:(
			&net.Dialer{
				Timeout: 30 * time.Second, //连接超时
				/*
				 Keep-Alive：使客户端到服务器端的连接持续有效，当出现对服务器的后继请求时，
				 Keep-Alive功能避免了建立或者重新建立连接。Web服务器，基本上都支持HTTP Keep-Alive。
				*/
				KeepAlive: 30 * time.Second, //探活时间
			}).DialContext,
			MaxIdleConns: 100, //支持最大空闲连接
			IdleConnTimeout: 90 * time.Second, //空闲超时连接
			TLSHandshakeTimeout: 90 * time.Second, //闲置连接在连接池的保存时间
			ExpectContinueTimeout: 1 * time.Second, //限制client在发送包含 Expect: 100-continue的header到收到继续发送body的response之间的时间等待。
	}
	client := &http.Client{
		Timeout: 50 * time.Second, //请求的超时时间
		Transport: transport, //连接池
	}
	resp,err := client.Get("http://127.0.0.1:1210/bye")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bsd,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bsd))
}
