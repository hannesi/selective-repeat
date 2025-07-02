package clientprotocol

import (
	"log"
	"math"
	"slices"
	"time"

	"github.com/hannesi/selective-repeat/internal/config"
	"github.com/hannesi/selective-repeat/internal/reliability"
	"github.com/hannesi/selective-repeat/internal/virtualsocket"
	"github.com/hannesi/selective-repeat/pkg/utils"
)

type SelectiveRepeatProtocolClient struct {
	socket    virtualsocket.VirtualSocket
	sequencer utils.Sequencer
	// packetQueue could be a struct of its own with methods
	packetQueue []reliability.ReliableDataTransferPacket
}

func NewSelectiveRepeatProtocolClient(socket virtualsocket.VirtualSocket) (SelectiveRepeatProtocolClient, error) {
	client := SelectiveRepeatProtocolClient{
		socket:      socket,
		sequencer:   utils.NewSequencer(config.DefaultConfig.GoBackNMaxSequence),
		packetQueue: []reliability.ReliableDataTransferPacket{},
	}

	helloAttempts := 0
	gotResponse := false
	for helloAttempts < config.DefaultConfig.HelloCountBeforeQuit && !gotResponse {
		gotResponse, _ = client.sendHello()
		helloAttempts++
	}

	if !gotResponse {
		return SelectiveRepeatProtocolClient{}, reliability.HelloError{}
	}

	return client, nil
}

func (client SelectiveRepeatProtocolClient) sendHello() (bool, error) {
	// TODO: replace the mystery constant below
	log.Println("Sending HELLO...")
	buffer := make([]byte, 5)

	client.socket.Send([]byte(config.DefaultConfig.HelloMessage))

	_, err := client.socket.Receive(buffer)
	if err != nil {
		return false, err
	}

	res := string(buffer[:])

	log.Printf("Received response to HELLO: %+v\n", res)

	return res == config.DefaultConfig.HelloMessage, err
}

func (client SelectiveRepeatProtocolClient) Send(data [][]byte) error {
	// form rdt packets from each byte array
	for _, payload := range data {
		rdtPacket := reliability.NewReliableDataTransferPacket(client.sequencer.Next(), payload)
		client.packetQueue = append(client.packetQueue, rdtPacket)
	}

	// send queued packets
	client.sendPacketQueue()

	return nil
}

func (client SelectiveRepeatProtocolClient) sendPacketQueue() {
	batch := client.makeBatch()
	ackChan := make(chan uint8)
	go client.listenForAcks(ackChan)
	client.sendBatch(batch)
	highestAckSeqReceived := <-ackChan

	// Remove packets up to and including the one with highest seq acked
	// The searched index is expected to be near the beginning of the slice, so the chosen method is fine.
	highestAckedIdx := slices.IndexFunc(client.packetQueue, func(p reliability.ReliableDataTransferPacket) bool {
		return p.Sequence == highestAckSeqReceived
	})

	if highestAckedIdx != -1 {
		client.packetQueue = slices.Delete(client.packetQueue, 0, highestAckedIdx+1)
	}

	if len(client.packetQueue) != 0 {
		client.sendPacketQueue()
	} else {
		log.Println("All packets sent! Shutting down...")
	}
}

func (client SelectiveRepeatProtocolClient) listenForAcks(ackChannel chan uint8) {
	log.Print("Listening for acks...")
	var highestAckSeqReceived uint8
	startTime := time.Now()
	for time.Now().Sub(startTime) < config.DefaultConfig.GoBackNAckCollectingTime {
		buffer := make([]byte, 4)
		_, err := client.socket.Receive(buffer)
		if err != nil {
			continue
		}
		ack, err := reliability.DeserializeAckBytes(buffer)
		if err != nil {
			continue
		}
		highestAckSeqReceived = ack.Sequence
	}
	ackChannel <- highestAckSeqReceived
}

func (client SelectiveRepeatProtocolClient) makeBatch() [][]byte {
	n := int(math.Min(float64(config.DefaultConfig.GoBackNWindowSize), float64(len(client.packetQueue))))
	serializedPackets := make([][]byte, n)
	for i := range n {
		data, err := client.packetQueue[i].Serialize()
		if err != nil {
			log.Printf("%v", err)
		}
		serializedPackets[i] = data
	}
	return serializedPackets
}

func (client SelectiveRepeatProtocolClient) sendBatch(batch [][]byte) {
	log.Printf("%s==== Sending batch ====%s", config.PositiveHighlightColour, config.ResetColour)
	for _, packet := range batch {
		client.socket.Send(packet)
	}
}
