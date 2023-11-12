package main

import (
	"machine"
)

type UARTDevice struct {
	uart *machine.UART
	buf  []byte
	prev uint16
}

func NewUARTDevice() *UARTDevice {
	uart := machine.DefaultUART
	tx := machine.UART_TX_PIN
	rx := machine.UART_RX_PIN
	uart.Configure(machine.UARTConfig{TX: tx, RX: rx})

	for uart.Buffered() > 0 {
		// flush
		uart.ReadByte()
	}

	return &UARTDevice{
		uart: uart,
		buf:  make([]byte, 0, 3),
	}
}

func (u *UARTDevice) Get() uint16 {
	buttons := uint16(0)

	for u.uart.Buffered() > 0 {
		data, _ := u.uart.ReadByte()
		u.buf = append(u.buf, data)
		if len(u.buf) == 3 {
			if u.buf[0] != 0xFF {
				u.buf[0], u.buf[1] = u.buf[1], u.buf[2]
				u.buf = u.buf[:2]
			} else {
				buttons = (uint16(u.buf[1]) << 8) + uint16(u.buf[2])
				u.buf = u.buf[:0]
				u.prev = buttons
				return buttons
			}
		} else {
		}
	}

	return u.prev
}
