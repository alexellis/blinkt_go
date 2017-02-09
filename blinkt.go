package blinkt

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/alexellis/rpi"
)

const DAT int = 23
const CLK int = 24

// pulse sends a pulse through the DAT/CLK pins
func pulse(pulses int) {
	rpi.DigitalWrite(rpi.GpioToPin(DAT), 0)
	for i := 0; i < pulses; i++ {
		rpi.DigitalWrite(rpi.GpioToPin(CLK), 1)
		rpi.DigitalWrite(rpi.GpioToPin(CLK), 0)
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

func writeByte(val int) {
	for i := 0; i < 8; i++ {
		// 0b10000000 = 128
		rpi.DigitalWrite(rpi.GpioToPin(DAT), val&128)
		rpi.DigitalWrite(rpi.GpioToPin(CLK), 1)
		val = val << 1
		rpi.DigitalWrite(rpi.GpioToPin(CLK), 0)
	}
}

// SetClearOnExit turns all pixels off on Control + C / os.Interrupt signal.
func (bl *Blinkt) SetClearOnExit(clearOnExit bool) {

	if clearOnExit {

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
}

// Delay maps to time.Sleep, for ms milliseconds
func Delay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// Clear sets all the pixels to off, you still have to call Show.
func (bl *Blinkt) Clear() {
	r := 0
	g := 0
	b := 0
	bl.SetAll(r, g, b)
}

// Show updates the LEDs with the values from SetPixel/Clear.
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

// SetAll sets all pixels to specified r, g, b colour. Show must be called to update the LEDs.
func (bl *Blinkt) SetAll(r int, g int, b int) {
	for i, _ := range bl.pixels {
		bl.SetPixel(i, r, g, b)
	}
}

// SetPixel sets an individual pixel to specified r, g, b colour. Show must be called to update the LEDs.
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

// Setup initializes GPIO via WiringPi base library.
func (bl *Blinkt) Setup() {
	rpi.WiringPiSetup()
	rpi.PinMode(rpi.GpioToPin(DAT), rpi.OUTPUT)
	rpi.PinMode(rpi.GpioToPin(CLK), rpi.OUTPUT)
}

// NewBlinkt creates a Blinkt to interact with. You must call "Setup()" immediately afterwards.
func NewBlinkt(brightness int) Blinkt {
	return Blinkt{
		pixels: initPixels(brightness),
	}
}

// Blinkt use the NewBlinkt function to initialize the pixels property.
type Blinkt struct {
	pixels [8][4]int
}

func init() {

}
