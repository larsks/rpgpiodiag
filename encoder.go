package main

import (
	"fmt"
	"machine"
)

var encoderTable = [16]int8{
	0, -1, 1, 0,
	1, 0, 0, -1,
	-1, 0, 0, 1,
	0, 1, -1, 0,
}

type RotaryEncoder struct {
	Name           string
	PinA           machine.Pin
	PinB           machine.Pin
	StepsPerDetent int
	prevState      uint8
	position       int
}

func NewRotaryEncoder(name string, pinA, pinB machine.Pin) *RotaryEncoder {
	pinA.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	pinB.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	return &RotaryEncoder{
		Name:           name,
		PinA:           pinA,
		PinB:           pinB,
		StepsPerDetent: 4,
		prevState:      readPinState(pinA, pinB),
	}
}

func readPinState(pinA, pinB machine.Pin) uint8 {
	var state uint8
	if pinA.Get() {
		state |= 2
	}
	if pinB.Get() {
		state |= 1
	}
	return state
}

func (e *RotaryEncoder) Update() {
	curr := readPinState(e.PinA, e.PinB)
	index := (e.prevState << 2) | curr
	delta := encoderTable[index]
	e.prevState = curr

	if delta == 0 {
		return
	}

	e.position += int(delta)
	if e.position >= e.StepsPerDetent {
		fmt.Printf("%s: CW\n", e.Name)
		e.position = 0
	} else if e.position <= -e.StepsPerDetent {
		fmt.Printf("%s: CCW\n", e.Name)
		e.position = 0
	}
}

func (e *RotaryEncoder) Pins() []machine.Pin {
	return []machine.Pin{e.PinA, e.PinB}
}
