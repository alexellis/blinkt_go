package sysfs

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/alexellis/blinkt_go/sysfs/gpio"
)

const DAT int = 23
const CLK int = 24

func writeByte(val int) {
	for i := 0; i < 8; i++ {
		// 0b10000000 = 128
		gpio.DigitalWrite(gpio.GpioToPin(DAT), val&128)
		gpio.DigitalWrite(gpio.GpioToPin(CLK), 1)
		val = val << 1
		gpio.DigitalWrite(gpio.GpioToPin(CLK), 0)
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
				gpio.Cleanup()
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
	for i := 0; i < 4; i++ {
		writeByte(0)
	}

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
	writeByte(255) // 0xff = 255
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
	gpio.Setup()
	gpio.PinMode(gpio.GpioToPin(DAT), gpio.OUTPUT)
	gpio.PinMode(gpio.GpioToPin(CLK), gpio.OUTPUT)
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

