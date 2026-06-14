package main

import (
	"image/color"
	"machine"

	"tinygo.org/x/drivers/ws2812"
)

type (
	ledFlasher struct {
		Pin      machine.Pin
		Color    color.RGBA
		MaxTicks int
		ticks    int
		on       bool
		ws       ws2812.Device
	}
)

func NewFlasher() *ledFlasher {
	return &ledFlasher{
		Pin:      machine.WS2812,
		MaxTicks: 100,
		Color:    color.RGBA{R: 255},
	}
}

func (f *ledFlasher) SetPin(pin machine.Pin) *ledFlasher {
	f.Pin = pin
	return f
}

func (f *ledFlasher) SetColor(color color.RGBA) *ledFlasher {
	f.Color = color
	return f
}

func (f *ledFlasher) SetTicks(ticks int) *ledFlasher {
	f.MaxTicks = ticks
	return f
}

func (f *ledFlasher) Build() *ledFlasher {
	f.Pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	f.ws = ws2812.New(f.Pin)
	return f
}

func (f *ledFlasher) Update() {
	f.ticks++
	if f.ticks >= f.MaxTicks {
		f.ticks = 0
		f.on = !f.on
		if f.on {
			f.ws.WriteColors([]color.RGBA{f.Color})
		} else {
			f.ws.WriteColors([]color.RGBA{color.RGBA{}})
		}
	}
}
