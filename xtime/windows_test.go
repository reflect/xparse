package xtime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWindowsLocationResolution(t *testing.T) {
	locations := []struct {
		IANAName    string
		WindowsName string
	}{
		{"Africa/Bujumbura", "South Africa Standard Time"},
		{"Etc/UTC", "UTC"},
		{"America/New_York", "Eastern Standard Time"},
	}

	for _, location := range locations {
		l, err := time.LoadLocation(location.IANAName)
		require.NoError(t, err)

		name, err := WindowsLocationString(l)
		require.Nil(t, err)
		require.Equal(t, location.WindowsName, name)
	}
}
