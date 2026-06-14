# rpgpiodiag -- diagnostic firmware for the Raspberry Pi Pico

This project implements raspberry pi pico firmware (specifically targeting the Waveshare RP2040-Zero) for testing input device connections. It will monitor a list of GPIO pins and print whenever the input value changes. It also has support for rotary encoders; it uses an interrupt-based mechanism to ensure robust decoding.
