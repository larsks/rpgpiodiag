package main

import (
	"fmt"
	"image/color"
	"machine"
	"time"
)

// We will monitor these GPIO pins for changes. They will be configured as
// inputs with pull-ups enabled.
var gpioPins = []machine.Pin{
	machine.GPIO0, machine.GPIO1, machine.GPIO2, machine.GPIO3,
	machine.GPIO4, machine.GPIO5, machine.GPIO6, machine.GPIO7,
	machine.GPIO8, machine.GPIO9, machine.GPIO10, machine.GPIO11,
	machine.GPIO12, machine.GPIO13, machine.GPIO14,
	machine.GPIO17, machine.GPIO18, machine.GPIO19,
	machine.GPIO20, machine.GPIO21, machine.GPIO22, machine.GPIO23,
	machine.GPIO26, machine.GPIO27, machine.GPIO28, machine.GPIO29,
}

// We will monitor these rotary encoders.
var encoders = []*RotaryEncoder{
	NewRotaryEncoder("ENC0", machine.GPIO3, machine.GPIO2),
	NewRotaryEncoder("ENC1", machine.GPIO6, machine.GPIO5),
}

var (
	ledRed = color.RGBA{R: 128}
)

// Produce optional detail information for a GPIO pin. For rotary encoders it
// prints the encoder and signal name.
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

// Print the state of  all monitored GPIO pins.
func printAllPins(prevState []bool) {
	fmt.Println("=== GPIO Pin State ===")
	for i, pin := range gpioPins {
		prevState[i] = pin.Get()
		fmt.Printf("GPIO%-2d = %d%s\n", pin, boolToInt(prevState[i]), pinLabel(pin))
	}
	fmt.Println("=== Monitoring for changes ===")
}

func main() {
	flasher := NewFlasher().SetColor(ledRed).Build()

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

	for {
		flasher.Update()

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
			// Skip if this pin is claimed by a rotary encoder
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
