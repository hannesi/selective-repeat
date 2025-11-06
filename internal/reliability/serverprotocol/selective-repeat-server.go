package serverprotocol

import (
	"log"
	"maps"
	"net"

	"github.com/hannesi/selective-repeat/internal/config"
	"github.com/hannesi/selective-repeat/internal/reliability"
)

type SelectiveRepeatProtocolServer struct {
	Socket *net.UDPConn
}

func NewSelectiveRepeatProtocolServer(socket *net.UDPConn) SelectiveRepeatProtocolServer {
	return SelectiveRepeatProtocolServer{
		Socket: socket,
	}
}

func (server SelectiveRepeatProtocolServer) Receive(msgChan chan []byte) error {
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
			server.sendAck(addr, packet.Sequence)
			if packet.Sequence == windowBase {
				msgChan <- packet.Payload
				// check packetBuffer for subsequent packets
				windowBase++
				for {
					message, exists := packetBuffer[windowBase]
					if !exists {
						break
					}
					msgChan <- message
					delete(packetBuffer, windowBase)
					windowBase++
				}
			} else {
				// if out of order, add to buffer
				packetBuffer[packet.Sequence] = packet.Payload
			}
		case false:
			log.Println("Reliability layer: Bit error detected!")
		}
	}
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
