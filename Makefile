TARGET = waveshare-rp2040-zero
OUTPUT = rpgpiodiag.uf2

.PHONY: build
build:
	tinygo build -target=$(TARGET) -o $(OUTPUT) .
