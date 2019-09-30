package main

import (
	"conn/udp"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("-----begin server------")
	u, err := udp.NewUDPServer(udp.UDP_SERVER_ADDR)
	if err != nil {
		logrus.Error("New UDPServer error:", err)
		os.Exit(-1)
	}
	u.Use(&udp.TestReceive{}) // you can use your struct ,need implement Spit
	u.Receive()
	fmt.Println("-----end server------")
}
