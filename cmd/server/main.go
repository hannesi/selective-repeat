package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/hannesi/selective-repeat/internal/config"
	"github.com/hannesi/selective-repeat/internal/reliability/serverprotocol"
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

	msgChan := make(chan []byte)
	reliabilityLayer := serverprotocol.NewSelectiveRepeatProtocolServer(socket, msgChan)

	go reliabilityLayer.Receive()

	for {
		msg := <-msgChan
		log.Printf("Received message: %s%s%s", config.PositiveHighlightColour, string(msg), config.ResetColour)
	}

}
