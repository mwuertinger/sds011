package sds011

import (
	"errors"
	"fmt"
	"io"
)

type Message struct {
	Command uint8   // command
	PM25    float64 // pm2.5 concentraction in μg/m^3
	PM10    float64 // pm10 concentraction in μg/m^3
	ID      uint16
}

func ReadMessage(r io.Reader) (*Message, error) {
	var buf [10]byte
	// search for message header
	// all messages start with 0xAA
	for {
		_, err := r.Read(buf[0:1])
		if err != nil {
			return nil, err
		}
		if buf[0] == 0xAA {
			break
		}
	}

	_, err := io.ReadFull(r, buf[1:])
	if err != nil {
		return nil, err
	}

	return parseMessage(buf)
}

func parseMessage(buf [10]byte) (*Message, error) {
	if buf[0] != 0xAA {
		return nil, fmt.Errorf("invalid message header: %x", buf[0])
	}
	if buf[9] != 0xAB {
		return nil, fmt.Errorf("invalid message tail: %x", buf[0])
	}

	var checksum byte
	for _, b := range buf[2:8] {
		checksum += b
	}
	if checksum != buf[8] {
		return nil, errors.New("checksum mismatch")
	}

	var pm25, pm10 uint16
	pm25 = uint16(buf[2])
	pm25 |= uint16(buf[3]) << 8
	pm10 = uint16(buf[4])
	pm10 |= uint16(buf[5]) << 8

	var msg Message
	msg.Command = buf[1]
	msg.PM25 = float64(pm25) / 10
	msg.PM10 = float64(pm10) / 10
	msg.ID = uint16(buf[6])
	msg.ID |= uint16(buf[7]) << 8

	return &msg, nil
}

func SetWorkingPeriod(w io.Writer, minutes uint8) error {
	var cmd [15]byte
	cmd[0] = 8
	cmd[1] = 1
	cmd[2] = minutes
	_, err := w.Write(cmd[:])
	return err
}

func sendCommand(w io.Writer, cmd [15]byte, deviceId uint16) error {
	var buf [19]byte
	buf[0] = 0xAA
	buf[1] = 0xB4
	for i := 0; i < 15; i++ {
		buf[i+2] = cmd[i]
	}
	// send to all device IDs
	buf[15] = 0xFF
	buf[16] = 0xFF
	// calculate checksum
	for _, b := range cmd {
		buf[17] += b
	}
	buf[18] = 0xAB

	_, err := w.Write(buf[:])
	return err
}
