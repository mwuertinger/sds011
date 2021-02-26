package sds011

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadMessage(t *testing.T) {
	data := []struct {
		testName string
		msg      [10]byte
		expected *Message
	}{
		{
			"standard reading",
			[10]byte{0xAA, 0xC0, 0xD4, 0x04, 0x3A, 0x0A, 0xA1, 0x60, 0x1D, 0xAB},
			&Message{
				Command: 0xC0,
				PM25:    123.6,
				PM10:    261.8,
				ID:      0x60A1,
			},
		},
		{
			"set working period reply",
			[10]byte{0xAA, 0xC5, 0x08, 0x01, 0x01, 0x00, 0xA1, 0x60, 0x0B, 0xAB},
			&Message{
				Command: 0xC5,
				PM25:    26.4,
				PM10:    0.1,
				ID:      0x60A1,
			},
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			out, err := ReadMessage(bytes.NewReader(d.msg[:]))
			assert.Nil(t, err)
			assert.Equal(t, d.expected, out)
		})
	}
}

func TestSendCommand(t *testing.T) {
	data := []struct {
		cmd         [15]byte
		expectedOut [19]byte
	}{
		{
			[15]byte{8, 1, 17, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF, 0xFF},
			[19]byte{0xAA, 0xB4, 8, 1, 17, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF, 0xFF, 0x18, 0xAB},
		},
	}

	for _, d := range data {
		var out bytes.Buffer
		sendCommand(&out, d.cmd)
		assert.Equal(t, d.expectedOut[:], out.Bytes(), "")
	}
}
