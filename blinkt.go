package blinkt

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	. "github.com/alexellis/rpi"
)

const DAT int = 23
const CLK int = 24

// pulse sends a pulse through the DAT/CLK pins
func pulse(pulses int) {
	DigitalWrite(GpioToPin(DAT), 0)
	for i := 0; i < pulses; i++ {
		DigitalWrite(GpioToPin(CLK), 1)
		DigitalWrite(GpioToPin(CLK), 0)
	}
}

// eof end of file or signal, from Python library
func eof() {
	pulse(36)
}

// sof start of file (name from Python library)
func sof() {
	pulse(32)
}

func (bl *Blinkt) Setup() {
	WiringPiSetup()
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

func (bl *Blinkt) InitStopper() {

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	fmt.Println("Press Control + C to stop")

	go func() {
		for range signalChan {
			bl.Clear()
			bl.Show()
			os.Exit(1)
		}
	}()
}

func (bl *Blinkt) Sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func (bl *Blinkt) Clear() {
	r := 0
	g := 0
	b := 0
	bl.SetAll(r, g, b)
}

func (bl *Blinkt) Show() {
	sof()
	for i, _ := range bl.pixels {
		brightness := bl.pixels[i][3]
		r := bl.pixels[i][0]
		g := bl.pixels[i][1]
		b := bl.pixels[i][2]

		// 0b11100000 (224)
		bitwise := 224
		writeByte(bitwise | brightness)
		writeByte(b)
		writeByte(g)
		writeByte(r)
	}
	eof()
}

func (bl *Blinkt) SetAll(r int, g int, b int) {
	for i, _ := range bl.pixels {
		bl.SetPixel(i, r, g, b)
	}
}

func (bl *Blinkt) SetPixel(p int, r int, g int, b int) {
	bl.pixels[p][0] = r
	bl.pixels[p][1] = g
	bl.pixels[p][2] = b
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

func NewBlinkt(brightness int) Blinkt {
	return Blinkt{
		pixels: initPixels(brightness),
	}
}

type Blinkt struct {
	pixels [8][4]int
}

func init() {

}
