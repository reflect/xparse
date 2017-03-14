package xtime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimezoneOffset(t *testing.T) {
	winter, _ := time.Parse("2006-01-02 15:04:05", "2016-01-01 05:00:00")
	summer, _ := time.Parse("2006-01-02 15:04:05", "2016-07-01 05:00:00")

	fixtures := []struct {
		Where          string
		At             time.Time
		Hours, Minutes int
	}{
		{"Africa/Bujumbura", winter, 2, 0},
		{"America/Curacao", winter, -4, 0},
		{"Asia/Pyongyang", winter, 8, 30},
		{"Canada/Pacific", winter, -8, 0},
		{"Canada/Pacific", summer, -7, 0},
	}

	for _, fixture := range fixtures {
		hours, mins, err := TimezoneOffsetAt(fixture.Where, fixture.At)

		assert.NoError(t, err)
		assert.Equal(t, fixture.Hours, hours)
		assert.Equal(t, fixture.Minutes, mins)
	}
}
