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
	scale    bool
}

func NewADCDevice(pin machine.Pin, num, min, max int, scale bool) *ADCDevice {
	adc := machine.ADC{Pin: pin}
	adc.Configure(machine.ADCConfig{})
	return &ADCDevice{
		buf:   make([]uint16, num),
		adc:   adc,
		min:   min,
		max:   max,
		scale: scale,
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
	ave = uint32(a.RawValue)

	ret := 32767 * (int(ave) - a.min) / (a.max - a.min)
	if ret < 0 {
		ret = 0
	}

	if a.scale {
		ret = scaleValue(ret)
	}

	if 32767 < ret {
		ret = 32767
	}
	a.Value = uint16(ret)
	return a.Value
}

func scaleValue(input int) int {
	// 最大値で正規化
	const maxUint16 = 0x7FFF
	normalized := float64(input) / float64(maxUint16)

	var scaled float64
	if normalized < 0.2 {
		// 0.0 ～ 0.2 の場合は 0.0 ～ 0.7 に線形にスケール
		scaled = normalized * 3.5
	} else {
		// 0.2 ～ 1.0 の場合は 0.7 ～ 1.0 に線形にスケール
		scaled = 0.7 + (normalized-0.2)*(0.3/0.8)
	}

	// 結果を int に変換して返す
	return int(scaled * float64(maxUint16))
}
