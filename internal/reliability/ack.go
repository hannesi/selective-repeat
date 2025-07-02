package reliability

import (
	"bytes"
	"io"
)

type AckData struct {
	Ack      string
	Sequence uint8
}

func NewAckPacket(msg string, sequence uint8) AckData {
	return AckData{
		Ack: msg,
		Sequence: sequence,
	}
}

// Serializes ack data to the following format: [sequence, ...msg]
func (a *AckData) Serialize() ([]byte, error) {
	buffer := new(bytes.Buffer)

	err := buffer.WriteByte(a.Sequence)
	if err != nil {
		return nil, err
	}

	_, err = buffer.Write([]byte(a.Ack))
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Deserializes ack data from the following format: [sequence, ...msg]
func DeserializeAckBytes(data []byte) (AckData, error) {
	buffer := bytes.NewReader(data)

	seq, err := buffer.ReadByte()
	if err != nil {
		return AckData{}, err
	}

	ack, err := io.ReadAll(buffer)
	if err != nil {
		return AckData{}, err
	}
	
	ackData := AckData{
		Ack: string(ack),
		Sequence: seq,
	}
	
	return ackData, nil
}
