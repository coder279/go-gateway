package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	url2 "net/url"
	"strings"
)

var addr = "127.0.0.1:2002"
func main() {
	rs1 := "http://127.0.0.1:2003/base"
	url,err := url2.Parse(rs1) //url对象
	if err != nil {
		log.Fatal(err)
		return
	}
	proxy := NewSingleHostReverseProxy(url)
	log.Println("starting httpserver at" + addr)
	//监听这个地址 调用反向代理的方法
	log.Fatal(http.ListenAndServe(addr,proxy))
}
func NewSingleHostReverseProxy(target *url2.URL) *httputil.ReverseProxy {
	//http://127.0.0.1?name=123
	//RawQuery: name=123
	//Scheme: http
	//Host: http:127.0.0.1
	targetQuery := target.RawQuery //将目标的参数发过来
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme //目标协议是啥请求
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	modify := func(res *http.Response) error {
		if(res.StatusCode == http.StatusOK){
			return nil
		}
		oldPayLoad,err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		newPayLoader := []byte("hello "+string(oldPayLoad))
		res.Body = ioutil.NopCloser(bytes.NewBuffer(newPayLoader))
		res.ContentLength = int64(len(newPayLoader))
		res.Header.Set("Content-Length",fmt.Sprint(len(newPayLoader)))
		return nil
	}
	return &httputil.ReverseProxy{Director: director,ModifyResponse:modify}
}
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
