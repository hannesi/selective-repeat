package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/hannesi/go-back-n/internal/config"
	"github.com/hannesi/go-back-n/internal/reliability/serverprotocol"
)

func main() {
	fmt.Println("SERVER")

	addr := net.UDPAddr{
		IP:   net.ParseIP(config.DefaultConfig.IPAddrString),
		Port: config.DefaultConfig.ServerPort,
	}

	socket, err := net.ListenUDP("udp", &addr)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer socket.Close()

	reliabilityLayer := serverprotocol.NewGoBackNProtocolServer(socket)
	msgChan := make(chan string)

	go reliabilityLayer.Receive(msgChan)

	for {
		msg := <-msgChan
		log.Printf("Received message: %s%s%s", config.PositiveHighlightColour, msg, config.ResetColour)
	}

}
