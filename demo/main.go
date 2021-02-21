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

	for {
		msg, err := sds011.ReadMessage(s)
		if err == io.EOF {
			log.Fatalf("%v", err)
		}
		if err != nil {
			log.Printf("%v", err)
			continue
		}
		log.Printf("message: pm2.5=%.1f pm10=%.1f", msg.PM25, msg.PM10)
	}
}
