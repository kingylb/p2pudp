// UDPServer project main.go
package main

import (
	"p2pudp/server"
	"time"
)

const (
	ListenPort = "10000"
	)

func main() {

	udpfuns := server.UDPHandle{ListenPort}
	udpfuns.MakUDPListener()
	defer udpfuns.CloseUDPConn()
	udpfuns.ReciveData()

	for {
		time.Sleep(time.Second)
	}

}
