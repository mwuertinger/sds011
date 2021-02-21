package main

import (
	"io"
	"log"

	"github.com/mwuertinger/sds011"
	"github.com/tarm/serial"
)

func main() {
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal("open port: ", err)
	}
	defer s.Close()

	var buf [10]byte
	for {
		// search for message header
		for {
			_, err := s.Read(buf[0:1])
			if err != nil {
				log.Printf("read: %v", err)
				continue
			}
			if buf[0] == 0xAA {
				break
			}
		}

		_, err = io.ReadFull(s, buf[1:])
		if err != nil {
			log.Printf("read: %v", err)
			continue
		}
		msg, err := sds011.ParseMessage(buf)
		if err != nil {
			log.Printf("parse: %v", err)
			continue
		}
		log.Printf("message: pm2.5=%.1f pm10=%.1f", msg.PM25, msg.PM10)
	}
}
