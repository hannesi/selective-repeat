package main

import (
	"fmt"

	"github.com/hannesi/go-back-n/internal/reliability/clientprotocol"
	"github.com/hannesi/go-back-n/internal/virtualsocket"
)

var predefinedMessages = []string{"Alekhine", "Botvinnik", "Capablanca", "Ding", "Euwe", "Finegold", "Giri", "Houska", "Ivanchuk", "Jaenisch", "Karpov", "Löwenthal", "Muzychuk", "Naroditsky", "Ojanen", "Polugaevsky", "Qin", "Réti", "Shirov", "Tal", "Ushenina", "Vachier-Lagrave", "Williams", "Xie", "Yusupov", "Zaitsev"}

func main() {
	fmt.Println("Client")

	socket, err := virtualsocket.NewVirtualSocket()

	if err != nil {
		panic("Failed to create virtual socket.")
	}

	defer socket.Close()

	gbn, err := clientprotocol.NewGoBackNProtocolClient(socket)

	if err != nil {
		panic(fmt.Sprintf("Failed to start GBN protocol client: %v", err))
	}

	var messages [][]byte
	for _, msg := range predefinedMessages {
		messages = append(messages, []byte(msg))
	}

	gbn.Send(messages)
}

