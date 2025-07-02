package config

import "time"

type Config struct {
	GoBackNAckCollectingTime    time.Duration
	GoBackNMaxSequence          uint8
	GoBackNWindowSize           uint8
	HelloMessage                string
	HelloCountBeforeQuit        int
	IPAddrString                string
	ServerPort                  int
	VirtualSocketDelayRate      float64
	VirtualSocketDelay          time.Duration
	VirtualSocketDropRate       float64
	VirtualSocketErrorRate      float64
	ReliabilityLayerAckWaitTime time.Duration
	ServerPacketHandleTime      time.Duration
}

var DefaultConfig = Config{
	GoBackNAckCollectingTime:    1 * time.Second,
	GoBackNMaxSequence:          ^uint8(0),
	GoBackNWindowSize:           5,
	HelloMessage:                "HELLO",
	HelloCountBeforeQuit:        5,
	IPAddrString:                "127.0.0.1",
	ServerPort:                  42069,
	VirtualSocketDelayRate:      0.1,
	VirtualSocketDelay:          1500 * time.Millisecond,
	VirtualSocketDropRate:       0.1,
	VirtualSocketErrorRate:      0.1,
	ReliabilityLayerAckWaitTime: 240 * time.Millisecond,
	ServerPacketHandleTime:      200 * time.Millisecond,
}

const (
	ResetColour             = "\033[0m"
	PositiveHighlightColour = "\033[32m"
	NegativeHighlightColour = "\033[31m"
)
