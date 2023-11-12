package main

import (
	"context"
	"fmt"
	"log"
	"machine"
	"time"

	"tinygo.org/x/drivers/mcp2515"

	"github.com/SWITCHSCIENCE/ffb_steering_controller/control"
	"github.com/SWITCHSCIENCE/ffb_steering_controller/settings"
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
	s := settings.Get()
	s.NeutralAdjust = 20
	settings.Update(s)

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

	//js.SetNeutralAdjust(-30)

	machine.InitADC()
	//accel := NewADCDevice(machine.A1, 8, 0x1E00, 0x3480)
	accel := NewADCDevice(machine.A1, 8, 0x2A00, 0x3480)
	brake := NewADCDevice(machine.A0, 8, 0xE380, 0xF600)

	uartButtons := NewUARTDevice()

	go func() {
		tick := time.NewTicker(10 * time.Millisecond)
		cnt := 0
		for {
			select {
			case <-tick.C:
				// ここに自分の実装を書く
				js.SetAxis(2, int(accel.Get()))
				js.SetAxis(4, int(brake.Get()))

				btn := uartButtons.Get()
				js.SetButton(2, (btn&0x0008) > 0)
				js.SetButton(3, (btn&0x0004) > 0)
				js.SetButton(4, (btn&0x8000) > 0)
				js.SetButton(5, (btn&0x4000) > 0)
				js.SetButton(6, (btn&0x2000) > 0)
				js.SetButton(7, (btn&0x1000) > 0)
				js.SetButton(8, (btn&0x0800) > 0)
				js.SetButton(9, (btn&0x0400) > 0)
				js.SetButton(10, (btn&0x0200) > 0)
				js.SetButton(11, (btn&0x0100) > 0)
				js.SetButton(12, (btn&0x0080) > 0)
				js.SetButton(13, (btn&0x0040) > 0)
				js.SetButton(14, (btn&0x0020) > 0)
				js.SetButton(15, (btn&0x0010) > 0)

				if false && (cnt%50) == 0 {
					fmt.Printf("%04X %04X : %04X %04X\n",
						accel.RawValue, accel.Value,
						brake.RawValue, brake.Value,
					)
				}
				cnt++
			}
		}
	}()
	for {
		if err := js.Loop(ctx); err != nil {
			log.Print(err)
			time.Sleep(3 * time.Second)
		}
	}
}
