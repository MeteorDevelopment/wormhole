package channel

import (
	"wormhole/pkg/models/message"
	"wormhole/pkg/snowflake"
)

type Channel struct {
	Id       uint64
	Name     string
	Messages []*message.Message
}

func New(name string) *Channel {
	return &Channel{
		Id:       snowflake.NextId(),
		Name:     name,
		Messages: []*message.Message{},
	}
}

func Default() *Channel {
	return New("general")
}
