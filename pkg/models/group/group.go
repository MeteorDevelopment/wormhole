package group

import (
	"wormhole/pkg/models/account"
	"wormhole/pkg/models/channel"
	"wormhole/pkg/snowflake"
)

type Group struct {
	Id       uint64
	Name     string
	Owner    uint64
	Admins   []uint64
	Members  []uint64
	Channels []*channel.Channel
}

func New(name string, owner *account.Account, members []*account.Account) *Group {
	ids := make([]uint64, len(members))
	ids = append(ids, owner.Id)
	for i, m := range members {
		ids[i] = m.Id
	}

	return &Group{
		Id:       snowflake.NextId(),
		Name:     name,
		Owner:    owner.Id,
		Admins:   []uint64{owner.Id},
		Members:  ids,
		Channels: []*channel.Channel{channel.Default()},
	}
}
