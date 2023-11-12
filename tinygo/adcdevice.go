package main

import (
	"machine"
)

type ADCDevice struct {
	buf      []uint16
	idx      int
	adc      machine.ADC
	min      int
	max      int
	RawValue uint16
	Value    uint16
}

func NewADCDevice(pin machine.Pin, num, min, max int) *ADCDevice {
	adc := machine.ADC{Pin: pin}
	adc.Configure(machine.ADCConfig{})
	return &ADCDevice{
		buf: make([]uint16, num),
		adc: adc,
		min: min,
		max: max,
	}
}

func (a *ADCDevice) Get() uint16 {
	a.RawValue = a.adc.Get()
	a.buf[a.idx] = a.RawValue
	a.idx = (a.idx + 1) % len(a.buf)

	sum := uint32(0)
	for _, v := range a.buf {
		sum += uint32(v)
	}
	ave := (sum / uint32(len(a.buf)))

	ret := 32767 * (int(ave) - a.min) / (a.max - a.min)
	if ret < 0 {
		ret = 0
	}
	if 32767 < ret {
		ret = 32767
	}
	a.Value = uint16(ret)
	return a.Value
}
