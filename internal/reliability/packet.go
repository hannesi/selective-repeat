package reliability

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
)

// Represents a packet used for detecting bit errors.
type ReliableDataTransferPacket struct {
    Sequence uint8 // Sequence number of the packet.
    Payload []byte  // The payload of the packet.
    Checksum uint32 // The checksum calculated for the payload.
}

// Create a new rdt packet from provided payload. The checksum is calculated automagically.
func NewReliableDataTransferPacket(sequence uint8, payload []byte) ReliableDataTransferPacket {
    packet := ReliableDataTransferPacket{
        Sequence: sequence,
        Payload: payload,
    }
    packet.Checksum = packet.computeChecksum()
    return packet
}

func (p *ReliableDataTransferPacket) computeChecksum() uint32 {
    data := append([]byte{p.Sequence}, p.Payload...)
    return crc32.ChecksumIEEE(data)
}

// Returns true if the rdt packet's checksum matches the packet's payload's calculated checksum.
func (p ReliableDataTransferPacket) IsChecksumValid() bool {
    return p.Checksum == p.computeChecksum()
}

// Serialize a rdt packet into transferable form.
func (p ReliableDataTransferPacket) Serialize() ([]byte, error) {
    buffer := new(bytes.Buffer)

    err := binary.Write(buffer, binary.BigEndian, p.Checksum)
    if err != nil {
        return nil, err
    }

    err = buffer.WriteByte(p.Sequence)
    if err != nil {
        return nil, err
    }

    _, err = buffer.Write(p.Payload)
    if err != nil {
        return nil, err
    }

    return buffer.Bytes(), nil
}


// DeserializeReliableDataTransferPacket a byte array into a rdt packet.
func DeserializeReliableDataTransferPacket(data []byte) (ReliableDataTransferPacket, error) {
    buffer := bytes.NewReader(data)

    var checksum uint32

    err := binary.Read(buffer, binary.BigEndian, &checksum)
    if err != nil {
        return ReliableDataTransferPacket{}, err
    }

    sequence, err := buffer.ReadByte()
    if err != nil {
        return ReliableDataTransferPacket{}, err
    }

    payload := make([]byte, buffer.Len())
    _, err = buffer.Read(payload)
    if err != nil {
        return ReliableDataTransferPacket{}, err
    }

    return ReliableDataTransferPacket{
        Payload: payload,
        Sequence: sequence,
        Checksum: checksum,
    }, nil
}

