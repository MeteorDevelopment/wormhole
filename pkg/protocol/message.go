package protocol

import (
	"encoding/json"
	"fmt"
	"wormhole/pkg/database"
)

type Message struct {
	Sender  *database.Account `json:"sender,omitempty"`
	Content string            `json:"content"`
}

func (m *Message) String() string {
	if m.Sender == nil {
		return m.Content
	}

	return fmt.Sprintf("%s: %s", m.Sender.Username, m.Content)
}

func (m *Message) Encode() ([]byte, error) {
	encodedMsg, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return encodedMsg, nil
}

func DecodeMessage(encodedMsg []byte, client *Client) (*Message, error) {
	msg := &Message{Sender: client.account}
	err := json.Unmarshal(encodedMsg, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func SystemMessage(content string) *Message {
	return &Message{Content: content}
}
