package snowflake

import "testing"

func TestSnowflakeGeneration(t *testing.T) {
	ids := make(map[uint64]bool)

	for i := 0; i < 100; i++ {
		id := NextId()
		if ids[id] {
			t.Errorf("Duplicate Id generated")
		} else {
			ids[id] = true
		}
	}
}
