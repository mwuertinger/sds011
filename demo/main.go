package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	influxdb "github.com/influxdata/influxdb1-client/v2"
	"github.com/mwuertinger/sds011"
	"github.com/tarm/serial"
)

type sniffer struct {
	s io.ReadWriter
}

func hexFmt(buf []byte) string {
	var out bytes.Buffer
	for _, b := range buf {
		fmt.Fprintf(&out, "%02x, ", b)
	}
	return out.String()
}

func (s *sniffer) Read(buf []byte) (int, error) {
	n, err := s.s.Read(buf)
	if err != nil {
		log.Printf("Sniffer: Read: %v", err)
		return 0, err
	}
	log.Printf("Sniffer: Read: %v", hexFmt(buf[0:n]))
	return n, nil
}

func (s *sniffer) Write(buf []byte) (int, error) {
	n, err := s.s.Write(buf)
	if err != nil {
		log.Printf("Sniffer: Write: %v", err)
		return 0, err
	}
	log.Printf("Sniffer: Write: %v", hexFmt(buf[0:n]))
	return n, nil
}

func main() {
	influxClient, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		log.Fatalf("influx: %v", err)
	}

	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal("open port: ", err)
	}
	defer s.Close()

	port := &sniffer{s}

	err = sds011.SetWorkingPeriod(s, 5)
	if err != nil {
		log.Fatalf("SetWorkingPeriod: %v", err)
	}

	for {
		msg, err := sds011.ReadMessage(port)
		if err == io.EOF {
			log.Fatalf("%v", err)
		}
		if err != nil {
			log.Printf("%v", err)
			continue
		}
		log.Printf("message: cmd=%02x pm2.5=%.1f pm10=%.1f", msg.Command, msg.PM25, msg.PM10)

		if msg.Command == 0xC0 {
			err = sendToInflux(influxClient, msg)
			if err != nil {
				log.Printf("sendToInflux: %v", err)
			}
		}
	}
}

func sendToInflux(influxClient influxdb.Client, msg *sds011.Message) error {
	// Create a new point batch
	bp, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database:  "sensors",
		Precision: "s",
	})
	if err != nil {
		return err
	}
	tags := map[string]string{"location": "Office"}
	fields := map[string]interface{}{}
	fields["pm2_5"] = msg.PM25
	fields["pm10"] = msg.PM10

	pt, err := influxdb.NewPoint("measurements", tags, fields, time.Now())
	if err != nil {
		return err
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := influxClient.Write(bp); err != nil {
		return err
	}
	return nil
}
