package sds011

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadMessage(t *testing.T) {
	data := []struct {
		msg      [10]byte
		expected *Message
	}{
		{
			[10]byte{0xAA, 0xC0, 0xD4, 0x04, 0x3A, 0x0A, 0xA1, 0x60, 0x1D, 0xAB},
			&Message{
				Command: 0xC0,
				PM25:    123.6,
				PM10:    261.8,
				ID:      0x60A1,
			},
		},
	}

	for _, d := range data {
		out, _ := ReadMessage(bytes.NewReader(d.msg[:]))
		assert.Equal(t, d.expected, out)
	}
}

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
