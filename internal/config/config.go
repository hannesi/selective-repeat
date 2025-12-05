package config

import "time"

type Config struct {
	IPAddrString                         string
	ServerPort                           int
	VirtualSocketDelayRate               float64
	VirtualSocketDelay                   time.Duration
	VirtualSocketDropRate                float64
	VirtualSocketErrorRate               float64
	ReliabilityLayerAckWaitTime          time.Duration
	SelectiveRepeatPacketResendWaitTime  time.Duration
	SelectiveRepeatSendLoopInterval      time.Duration
	SelectiveRepeatWindowSize            int
	ServerTopLayerChannelProcessingDelay time.Duration
}

var DefaultConfig = Config{
	IPAddrString:                         "127.0.0.1",
	ServerPort:                           42069,
	VirtualSocketDelayRate:               0.1,
	VirtualSocketDelay:                   500 * time.Millisecond,
	VirtualSocketDropRate:                0.1,
	VirtualSocketErrorRate:               0.1,
	ReliabilityLayerAckWaitTime:          25 * time.Millisecond,
	SelectiveRepeatPacketResendWaitTime:  300 * time.Millisecond,
	SelectiveRepeatSendLoopInterval:      100 * time.Millisecond,
	SelectiveRepeatWindowSize:            4,
	ServerTopLayerChannelProcessingDelay: 0 * time.Millisecond,
}

const (
	ResetColour             = "\033[0m"
	PositiveHighlightColour = "\033[32m"
	NegativeHighlightColour = "\033[31m"
)
