package reliability

import (
	"bytes"
	"fmt"

	"github.com/hannesi/go-back-n/internal/config"
)

type HelloError struct {
}

func (e HelloError) Error() string {
	return fmt.Sprint("No hello response from server.")
}

// TODO: obsolete

type HelloResponse struct {
	ExpectedSequence uint8
	MaxSequence      uint8
	WindowSize       uint8
}

func NewHelloResponse(expectedSequence uint8, maxSequence uint8, windowSize uint8) HelloResponse {
	return HelloResponse{
		ExpectedSequence: expectedSequence,
		MaxSequence:      maxSequence,
		WindowSize:       windowSize,
	}
}

func (res HelloResponse) Serialize() ([]byte, error) {
	buffer := new(bytes.Buffer)

	err := buffer.WriteByte(res.ExpectedSequence)
	if err != nil {
		return nil, err
	}

	err = buffer.WriteByte(res.MaxSequence)
	if err != nil {
		return nil, err
	}

	err = buffer.WriteByte(res.WindowSize)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func DeserializeHelloResponse(data []byte) (HelloResponse, error) {
	buffer := bytes.NewReader(data)

	currentSequence, err := buffer.ReadByte()
	if err != nil {
		return HelloResponse{}, err
	}

	maxSequence, err := buffer.ReadByte()
	if err != nil {
		return HelloResponse{}, err
	}

	windowSize, err := buffer.ReadByte()
	if err != nil {
		return HelloResponse{}, err
	}

	return HelloResponse{
		ExpectedSequence: currentSequence,
		MaxSequence:      maxSequence,
		WindowSize:       windowSize,
	}, nil
}

func IsHelloMessage(data []byte) bool {
	return string(data) == config.DefaultConfig.HelloMessage
}
