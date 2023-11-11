package main

import (
	"context"
	"log"
	"machine"
	"time"

	"tinygo.org/x/drivers/mcp2515"

	"github.com/SWITCHSCIENCE/ffb_steering_controller/control"
)

const (
	// for feather-rp2040 pins
	LED1      machine.Pin = machine.NoPin
	LED2      machine.Pin = machine.NoPin
	LED3      machine.Pin = machine.NoPin
	SW1       machine.Pin = machine.NoPin
	SW2       machine.Pin = machine.NoPin
	SW3       machine.Pin = machine.NoPin
	CAN_INT   machine.Pin = machine.NoPin
	CAN_RESET machine.Pin = machine.NoPin
	CAN_SCK   machine.Pin = 18
	CAN_TX    machine.Pin = 19
	CAN_RX    machine.Pin = 20
	CAN_CS    machine.Pin = 7
)

var (
	spi = machine.SPI0
)

func init() {
	LED1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	LED2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	LED3.Configure(machine.PinConfig{Mode: machine.PinOutput})
	LED1.High()
	LED2.High()
	LED3.High()
	SW1.Configure(machine.PinConfig{Mode: machine.PinInput})
	SW2.Configure(machine.PinConfig{Mode: machine.PinInput})
	SW3.Configure(machine.PinConfig{Mode: machine.PinInput})
	CAN_INT.Configure(machine.PinConfig{Mode: machine.PinInput})
	CAN_RESET.Configure(machine.PinConfig{Mode: machine.PinOutput})
	CAN_RESET.Low()
	time.Sleep(10 * time.Millisecond)
	CAN_RESET.High()
	time.Sleep(10 * time.Millisecond)
}

func main() {
	LED1.Low()
	log.SetFlags(log.Lmicroseconds)
	if err := spi.Configure(
		machine.SPIConfig{
			Frequency: 500000,
			SCK:       CAN_SCK,
			SDO:       CAN_TX,
			SDI:       CAN_RX,
			Mode:      0,
		},
	); err != nil {
		log.Print(err)
	}
	can := mcp2515.New(spi, CAN_CS)
	can.Configure()
	if err := can.Begin(mcp2515.CAN500kBps, mcp2515.Clock8MHz); err != nil {
		log.Fatal(err)
	}
	js := control.NewWheel(can)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		if err := js.Loop(ctx); err != nil {
			log.Print(err)
			time.Sleep(3 * time.Second)
		}
	}
}
