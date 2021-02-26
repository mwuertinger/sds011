package sds011

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendCommand(t *testing.T) {
	data := []struct {
		cmd         [15]byte
		deviceId    uint16
		expectedOut [19]byte
	}{
		{
			[15]byte{8, 1, 17, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			0xFFFF,
			[19]byte{0xAA, 0xB4, 8, 1, 17, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF, 0xFF, 0x1A, 0xAB},
		},
	}

	for _, d := range data {
		var out bytes.Buffer
		sendCommand(&out, d.cmd, d.deviceId)
		assert.Equal(t, d.expectedOut[:], out.Bytes(), "")
	}
}
