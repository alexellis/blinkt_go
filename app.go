package main

import (
	. "github.com/alexellis/rpi"
)

const DAT int = 23
const CLK int = 24

func show(pixels [8][4]int) {
	sof()
	for i, _ := range pixels {
		brightness := pixels[i][3]
		r := pixels[i][0]
		g := pixels[i][1]
		b := pixels[i][2]

		// 0b11100000 (224)
		bitwise := 224
		writeByte(bitwise | brightness)
		writeByte(b)
		writeByte(g)
		writeByte(r)
	}
	eof()
}

func initPixels(brightness int) [8][4]int {
	var pixels [8][4]int
	for i, _ := range pixels {
		pixels[i][0] = 0
		pixels[i][1] = 0
		pixels[i][2] = 0
		pixels[i][3] = brightness
	}
	return pixels
}

func pulse(pulses int) {
	DigitalWrite(GpioToPin(DAT), 0)
	for i := 0; i < pulses; i++ {
		DigitalWrite(GpioToPin(CLK), 1)
		DigitalWrite(GpioToPin(CLK), 0)
	}
}

func eof() {
	pulse(36)
}

func sof() {
	pulse(32)
}

func setup() {
	PinMode(GpioToPin(DAT), OUTPUT)
	PinMode(GpioToPin(CLK), OUTPUT)
}

func writeByte(val int) {
	for i := 0; i < 8; i++ {
		// 0b10000000 = 128
		DigitalWrite(GpioToPin(DAT), val&128)
		DigitalWrite(GpioToPin(CLK), 1)
		val = val << 1
		DigitalWrite(GpioToPin(CLK), 0)
	}
}

func Clear(pixels *[8][4]int) {
	r := 0
	g := 0
	b := 0
	for i := 0; i < 8; i++ {
		SetPixel(pixels, i, r, g, b)
	}
}

func SetPixel(pixels *[8][4]int, p int, r int, g int, b int) {
	pixels[p][0] = r
	pixels[p][1] = g
	pixels[p][2] = b
}

func main() {
	WiringPiSetup()
	pixels := initPixels(50)
	setup()

	Delay(100)

	r := 255
	g := 0
	b := 0
	for pixel := 0; pixel < 8; pixel++ {
		SetPixel(&pixels, pixel, r, g, b)
		show(pixels)
		Delay(100)
	}

	Delay(1000)
	Clear(&pixels)
	show(pixels)
}
