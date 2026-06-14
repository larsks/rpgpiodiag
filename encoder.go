package main

import (
	"fmt"
	"runtime/interrupt"
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
	enc := &RotaryEncoder{
		Name:           name,
		PinA:           pinA,
		PinB:           pinB,
		StepsPerDetent: 4,
		prevState:      readPinState(pinA, pinB),
	}
	handler := func(machine.Pin) {
		enc.handleInterrupt()
	}
	pinA.SetInterrupt(machine.PinToggle, handler)
	pinB.SetInterrupt(machine.PinToggle, handler)
	return enc
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

func (e *RotaryEncoder) handleInterrupt() {
	curr := readPinState(e.PinA, e.PinB)
	index := (e.prevState << 2) | curr
	delta := encoderTable[index]
	e.prevState = curr
	e.position += int(delta)
}

func (e *RotaryEncoder) Update() {
	state := interrupt.Disable()
	pos := e.position
	if pos >= e.StepsPerDetent {
		e.position -= e.StepsPerDetent
		interrupt.Restore(state)
		fmt.Printf("%s: CW\n", e.Name)
	} else if pos <= -e.StepsPerDetent {
		e.position += e.StepsPerDetent
		interrupt.Restore(state)
		fmt.Printf("%s: CCW\n", e.Name)
	} else {
		interrupt.Restore(state)
	}
}

func (e *RotaryEncoder) Pins() []machine.Pin {
	return []machine.Pin{e.PinA, e.PinB}
}
