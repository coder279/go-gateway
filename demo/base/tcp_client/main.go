package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	doSend()
	fmt.Println("dosend over")
	return
}

func doSend(){
	conn,err := net.Dial("tcp","localhost:9090")
	if err != nil {
		fmt.Printf("connect failed,err: %v\n",err.Error())
		return
	}
	defer conn.Close()
	//读取命令行
	inputReader := bufio.NewReader(os.Stdin)
	for{
		//一直读到\n
		input,err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Printf("read from console failed,err %v\n",err)
			break
		}
		//4.读到Q停止
		PauseInput := strings.TrimSpace(input)
		if PauseInput == "Q" {
			break
		}
		//5.回复服务器信息
		_,err = conn.Write([]byte(PauseInput))
		if err != nil {
			fmt.Printf("write failed,err %v\n",err)
			break
		}

	}
}