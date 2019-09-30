package main

import (
	"conn/udp"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("-----begin server------")
	u, err := udp.NewUDPServer(udp.UDP_SERVER_ADDR)
	if err != nil {
		log.Fatal("New UDPServer error:", err)
		os.Exit(-1)
	}
	u.Use(&udp.TestReceive{}) // you can use your struct ,need implement Spit
	u.Receive()
	fmt.Println("-----end server------")
}
