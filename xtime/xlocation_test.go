package xtime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type tzOffset struct {
	Hours   int
	Minutes int
}

func TestTimeZoneOffset(t *testing.T) {
	fixtures := map[string]tzOffset{
		"Africa/Bujumbura": {2, 0},
		"America/Curacao":  {-4, 0},
		"Asia/Pyongyang":   {8, 30},
		"Canada/Pacific":   {-8, 0},
	}

	for tz, fixture := range fixtures {
		hours, mins, err := TimezoneOffset(tz)

		assert.NoError(t, err)
		assert.Equal(t, fixture.Hours, hours)
		assert.Equal(t, fixture.Minutes, mins)
	}
}
