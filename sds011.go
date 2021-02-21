package sds011

import (
	"errors"
	"fmt"
	"io"
)

type Message struct {
	Commander uint8   // commander number
	PM25      float64 // pm2.5 concentraction in μg/m^3
	PM10      float64 // pm10 concentraction in μg/m^3
	ID        uint16
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
	msg.Commander = buf[1]
	msg.PM25 = float64(pm25) / 10
	msg.PM10 = float64(pm10) / 10
	msg.ID = uint16(buf[6])
	msg.ID |= uint16(buf[7]) << 8

	return &msg, nil
}
