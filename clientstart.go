package main

import (
	"time"
	"p2pudp/client"
	"p2pudp/web"
)

const (
	ListenPort = "1000"
	ServerIP ="10.1.193.147:10000"
	Webport = "9000"
)

func main() {

	clientObject :=  client.ClientHandle{ListenPort,ServerIP,""}
	clientObject.MakUDPListener()

	//开启外网服api 接口
	go func(){
		web.BuildServerWeb(Webport)
	}()

	clientObject.ClientUdpReciveData()
	clientObject.ConnectServer()
	defer clientObject.CloseClientUDPconn()

	for {
		time.Sleep(time.Second)
	}
}
