package protocol

import "encoding/json"

type PacketType int

const (
	Authenticate PacketType = iota
	Message
)

type Packet struct {
	Version int             `json:"version"`
	Type    PacketType      `json:"type"`
	Data    json.RawMessage `json:"data"`
}

func Outbound(pType PacketType, data []byte) *Packet {
	return &Packet{
		Version: Version,
		Type:    pType,
		Data:    data,
	}
}

func (p *Packet) Encode() ([]byte, error) {
	return json.Marshal(p)
}

func Decode(encoded []byte) (*Packet, error) {
	var p Packet
	err := json.Unmarshal(encoded, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
