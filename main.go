package main

import (
	"fmt"
	"machine"
	"time"
)

var gpioPins = []machine.Pin{
	machine.GPIO0, machine.GPIO1, machine.GPIO2, machine.GPIO3,
	machine.GPIO4, machine.GPIO5, machine.GPIO6, machine.GPIO7,
	machine.GPIO8, machine.GPIO9, machine.GPIO10, machine.GPIO11,
	machine.GPIO12, machine.GPIO13, machine.GPIO14,
	machine.GPIO16, machine.GPIO17, machine.GPIO18, machine.GPIO19,
	machine.GPIO20, machine.GPIO21, machine.GPIO22, machine.GPIO23,
	machine.GPIO26, machine.GPIO27, machine.GPIO28, machine.GPIO29,
}

func pinLabel(pin machine.Pin) string {
	for _, enc := range encoders {
		if pin == enc.PinA {
			return " " + enc.Name + "-A"
		}
		if pin == enc.PinB {
			return " " + enc.Name + "-B"
		}
	}
	return ""
}

func printAllPins(prevState []bool) {
	fmt.Println("=== GPIO Pin State ===")
	for i, pin := range gpioPins {
		prevState[i] = pin.Get()
		fmt.Printf("GPIO%-2d = %d%s\n", pin, boolToInt(prevState[i]), pinLabel(pin))
	}
	fmt.Println("=== Monitoring for changes ===")
}

var encoders = []*RotaryEncoder{
	NewRotaryEncoder("ENC0", machine.GPIO2, machine.GPIO3),
	NewRotaryEncoder("ENC1", machine.GPIO5, machine.GPIO6),
}

func main() {
	encoderPins := make(map[machine.Pin]bool)
	for _, enc := range encoders {
		for _, p := range enc.Pins() {
			encoderPins[p] = true
		}
	}

	for _, pin := range gpioPins {
		if !encoderPins[pin] {
			pin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
		}
	}

	prevState := make([]bool, len(gpioPins))
	printAllPins(prevState)

	for {
		for machine.Serial.Buffered() > 0 {
			b, err := machine.Serial.ReadByte()
			if err != nil {
				break
			}
			if b == '\r' || b == '\n' {
				printAllPins(prevState)
			}
		}

		for _, enc := range encoders {
			enc.Update()
		}

		for i, pin := range gpioPins {
			if encoderPins[pin] {
				continue
			}
			val := pin.Get()
			if val != prevState[i] {
				fmt.Printf("GPIO%-2d -> %d\n", pin, boolToInt(val))
				prevState[i] = val
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
