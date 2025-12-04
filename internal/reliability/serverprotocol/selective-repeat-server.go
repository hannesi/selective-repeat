package serverprotocol

import (
	"log"
	"net"
	"time"

	"github.com/hannesi/selective-repeat/internal/config"
	"github.com/hannesi/selective-repeat/internal/reliability"
)

type SelectiveRepeatProtocolServer struct {
	Socket            *net.UDPConn
	ChannelToTopLayer chan []byte
}

func NewSelectiveRepeatProtocolServer(socket *net.UDPConn, messageChannel chan []byte) SelectiveRepeatProtocolServer {
	return SelectiveRepeatProtocolServer{
		Socket:            socket,
		ChannelToTopLayer: messageChannel,
	}
}

func (server SelectiveRepeatProtocolServer) Receive() error {
	buffer := make([]byte, 1024)
	packetBuffer := map[uint8][]byte{}
	windowBase := uint8(0)
	for {
		n, addr, err := server.Socket.ReadFromUDP(buffer)
		if err != nil {
			continue
		}

		packet, err := reliability.DeserializeReliableDataTransferPacket(buffer[:n])

		switch packet.IsChecksumValid() {
		case true:
			if packet.Sequence == windowBase {
				server.sendAck(addr, packet.Sequence)
				server.sendToTopLayer(packet.Payload)
				// check packetBuffer for subsequent packets
				windowBase++
				for {
					message, exists := packetBuffer[windowBase]
					if !exists {
						break
					}
					server.sendToTopLayer(message)
					delete(packetBuffer, windowBase)
					windowBase++
				}
			} else if packet.Sequence < windowBase {
				log.Println("Duplicate out of order packet, not buffering.")
			} else {
				// if out of order, add to buffer
				server.sendAck(addr, packet.Sequence)
				packetBuffer[packet.Sequence] = packet.Payload
			}
		case false:
			log.Println("Reliability layer: Bit error detected!")
		}
	}
}

func (server SelectiveRepeatProtocolServer) sendToTopLayer(msg []byte) {
	time.Sleep(config.DefaultConfig.ServerTopLayerChannelProcessingDelay)
	server.ChannelToTopLayer <- msg
}

func (server SelectiveRepeatProtocolServer) sendAck(dest *net.UDPAddr, seq uint8) {
	ack := reliability.NewAckPacket("ACK", seq)
	serializedAck, err := ack.Serialize()
	if err != nil {
		log.Println("Server failed to serialize an ack packet. Trying again.")
		server.sendAck(dest, seq)
	}
	server.Socket.WriteToUDP(serializedAck, dest)
}
