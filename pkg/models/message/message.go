package message

import (
	"encoding/json"
)

type Message struct {
	Id        uint64
	GroupId   uint64
	ChannelId uint64
	SenderId  uint64
	Content   string
}

func (m *Message) Encode() ([]byte, error) {
	encodedMsg, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return encodedMsg, nil
}

func Decode(encodedMsg []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(encodedMsg, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func System(content string) *Message {
	return &Message{Content: content}
}
