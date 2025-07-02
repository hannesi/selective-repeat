package serverprotocol

import (
	"log"
	"net"
	"time"

	"github.com/hannesi/go-back-n/internal/config"
	"github.com/hannesi/go-back-n/internal/reliability"
)

type GoBackNProtocolServer struct {
	Socket    *net.UDPConn
	lastOkAck uint8
}

func NewGoBackNProtocolServer(socket *net.UDPConn) GoBackNProtocolServer {
	return GoBackNProtocolServer{
		Socket:    socket,
		lastOkAck: config.DefaultConfig.GoBackNMaxSequence,
	}
}

func (server GoBackNProtocolServer) Receive(msgChan chan string) error {
	buffer := make([]byte, 1024)
	n, addr, err := server.Socket.ReadFromUDP(buffer)
	if err != nil {
		return server.Receive(msgChan)
	}

	if reliability.IsHelloMessage(buffer[:n]) {
		// reset sequencer on HELLO
		server.lastOkAck = config.DefaultConfig.GoBackNMaxSequence
		log.Println("Received HELLO message, answering HELLO")
		server.sendHelloResponse(addr)
		return server.Receive(msgChan)
	}

	packet, err := reliability.DeserializeReliableDataTransferPacket(buffer[:n])
	time.Sleep(config.DefaultConfig.ServerPacketHandleTime)

	if !packet.IsChecksumValid() {
		log.Println("Bit error detected!")
		server.sendAck(addr)
		return server.Receive(msgChan)
	}

	if packet.Sequence != server.lastOkAck+1 {
		log.Printf("Unexpected sequence number! Expected %s%d%s, received %s%d%s",
			config.NegativeHighlightColour, server.lastOkAck+1, config.ResetColour,
			config.NegativeHighlightColour, packet.Sequence, config.ResetColour)
		server.sendAck(addr)
		return server.Receive(msgChan)
	}

	// if packet is ok
	server.lastOkAck = packet.Sequence
	server.sendAck(addr)

	msgChan <- string(packet.Payload)

	return server.Receive(msgChan)
}

func (server GoBackNProtocolServer) sendAck(dest *net.UDPAddr) {
	ack := reliability.NewAckPacket("ACK", server.lastOkAck)
	serializedAck, err := ack.Serialize()
	if err != nil {
		log.Println("Server failed to serialize an ack packet. Trying again.")
		server.sendAck(dest)
	}
	server.Socket.WriteToUDP(serializedAck, dest)
}

func (server GoBackNProtocolServer) sendHelloResponse(addr *net.UDPAddr) {
	server.Socket.WriteToUDP([]byte(config.DefaultConfig.HelloMessage), addr)
}
