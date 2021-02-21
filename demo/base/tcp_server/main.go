package main

import (
	"fmt"
	"net"
)

func main() {
	quit := make(chan int)
	listener,err := net.Listen("tcp","0.0.0.0:9090")
	if err != nil {
		fmt.Println("listen fail,err %v\n","结束连接")
		return
	}
	for{
		conn,err := listener.Accept()
		if err != nil {
			fmt.Printf("accept fail err:%v\n",err)
			continue
		}
		go progress(conn,quit)
		<-quit
		return


	}
}

func progress(conn net.Conn,signal chan<-int){
	for {
		var buf [128]byte
		n,err := conn.Read(buf[:])
		if err != nil {
			fmt.Printf("read from connect failed,err:%v\n","----end----")
			break
		}
		str := string(buf[:n])
		fmt.Printf("receive from client,data %v\n",str)
	}
	signal <- 1
	defer conn.Close()
}
