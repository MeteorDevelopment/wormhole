package snowflake

import (
	"github.com/sony/sonyflake"
	"strconv"
)

var snowflake = sonyflake.NewSonyflake(sonyflake.Settings{})

func NextId() uint64 {
	id, _ := snowflake.NextID()
	return id
}

func IdString(id uint64) string {
	return strconv.FormatUint(id, 10)
}
