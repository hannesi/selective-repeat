package clientprotocol

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/hannesi/selective-repeat/internal/config"
	"github.com/hannesi/selective-repeat/internal/reliability"
	"github.com/hannesi/selective-repeat/internal/virtualsocket"
)

type SelectiveRepeatProtocolClient struct {
	socket virtualsocket.VirtualSocket
}

type SelectiveRepeatPacket struct {
	SerializedRDTPacket []byte
	LatestSendTime      time.Time
	Acked               bool
}

func NewSelectiveRepeatProtocolClient(socket virtualsocket.VirtualSocket) (SelectiveRepeatProtocolClient, error) {
	client := SelectiveRepeatProtocolClient{
		socket: socket,
	}

	return client, nil
}

func (client SelectiveRepeatProtocolClient) Send(data [][]byte) error {
	packets := []SelectiveRepeatPacket{}
	for i, payload := range data {
		serializedRdtPacket, err := reliability.NewReliableDataTransferPacket(uint8(i), payload).Serialize()
		if err != nil {
			fmt.Println(err)
		}
		srPacket := SelectiveRepeatPacket{SerializedRDTPacket: serializedRdtPacket}
		packets = append(packets, srPacket)
	}

	// start ack listening thread
	ctx, killListenForAcks := context.WithCancel(context.Background())
	ackChannel := make(chan int)
	go client.ListenForAcks(ackChannel, ctx)

	windowBase := 0
	ackBuffer := []int{}

	for windowBase < len(packets) {
		fmt.Printf("%sStart of loop%s\n", config.NegativeHighlightColour, config.ResetColour)
		// handle received acks
	AckChannelLoop:
		for {
			select {
			case seq := <-ackChannel:
				fmt.Printf("Handled ack %d\n", seq)
				packets[seq].Acked = true
				if seq == windowBase {
					windowBase = seq + 1
				AckBufferLoop:
					for {
						idx := slices.Index(ackBuffer, windowBase)
						switch idx {
						case -1:
							break AckBufferLoop
						default:
							windowBase++
							ackBuffer = slices.Delete(ackBuffer, idx, idx+1)
						}
					}
				} else if seq > windowBase {
					ackBuffer = append(ackBuffer, seq)
				}
			default:
				break AckChannelLoop
			}
		}
		fmt.Printf("Window Base: %d, Buffered acks: %v\n", windowBase, ackBuffer)

		// send batch of packets
		for i := windowBase; i < windowBase+config.DefaultConfig.SelectiveRepeatWindowSize; i++ {
			if i < len(packets) && !packets[i].Acked && time.Since(packets[i].LatestSendTime) > config.DefaultConfig.SelectiveRepeatPacketResendWaitTime {
				client.socket.Send(packets[i].SerializedRDTPacket)
				packets[i].LatestSendTime = time.Now()
				fmt.Printf("Sending packet %d\n", i)
			}
		}

		time.Sleep(config.DefaultConfig.SelectiveRepeatSendLoopInterval)
	}

	killListenForAcks()
	return nil
}

func (client SelectiveRepeatProtocolClient) ListenForAcks(ackChannel chan int, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			buffer := make([]byte, 4)
			_, err := client.socket.Receive(buffer)
			if err != nil {
				// fmt.Println(err)
				break
			}

			ackPacket, err := reliability.DeserializeAckBytes(buffer)
			if err != nil {
				fmt.Println(err)
				break
			}

			// fmt.Printf("Received ack: %v\n", ackPacket)

			ackChannel <- int(ackPacket.Sequence)
		}
	}
}
