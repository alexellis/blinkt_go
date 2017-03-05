package sysfs

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/alexellis/blinkt_go/sysfs/gpio"
)

const DAT int = 23
const CLK int = 24

const redIndex int = 0
const greenIndex int = 1
const blueIndex int = 2
const brightnessIndex int = 3

// default raw brightness.  Not to be used user-side
const defaultBrightnessInt int = 15

//upper and lower bounds for user specified brightness
const minBrightness float64 = 0.0
const maxBrightness float64 = 1.0

func writeByte(val int) {
	for i := 0; i < 8; i++ {
		// 0b10000000 = 128
		gpio.DigitalWrite(gpio.GpioToPin(DAT), val&128)
		gpio.DigitalWrite(gpio.GpioToPin(CLK), 1)
		val = val << 1
		gpio.DigitalWrite(gpio.GpioToPin(CLK), 0)
	}
}

func convertBrightnessToInt(brightness float64) int {

	if !inRangeFloat(minBrightness, brightness, maxBrightness) {
		log.Fatalf("Supplied brightness was %#v - value should be between: %#v and %#v", brightness, minBrightness, maxBrightness)
	}

	return int(brightness * 31.0)

}

func inRangeFloat(minVal float64, testVal float64, maxVal float64) bool {

	return (testVal >= minVal) && (testVal <= maxVal)
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

	for p, _ := range bl.pixels {
		brightness := bl.pixels[p][brightnessIndex]
		r := bl.pixels[p][redIndex]
		g := bl.pixels[p][greenIndex]
		b := bl.pixels[p][blueIndex]

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
func (bl *Blinkt) SetAll(r int, g int, b int) *Blinkt {

	for p, _ := range bl.pixels {
		bl.SetPixel(p, r, g, b)
	}

	return bl
}

// SetPixel sets an individual pixel to specified r, g, b colour. Show must be called to update the LEDs.
func (bl *Blinkt) SetPixel(p int, r int, g int, b int) *Blinkt {

	bl.pixels[p][redIndex] = r
	bl.pixels[p][greenIndex] = g
	bl.pixels[p][blueIndex] = b

	return bl

}

// SetBrightness sets the brightness of all pixels. Brightness supplied should be between: 0.0 to 1.0
func (bl *Blinkt) SetBrightness(brightness float64) *Blinkt {

	brightnessInt := convertBrightnessToInt(brightness)

	for p, _ := range bl.pixels {
		bl.pixels[p][brightnessIndex] = brightnessInt
	}

	return bl
}

// SetPixelBrightness sets the brightness of pixel p. Brightness supplied should be between: 0.0 to 1.0
func (bl *Blinkt) SetPixelBrightness(p int, brightness float64) *Blinkt {

	brightnessInt := convertBrightnessToInt(brightness)
	bl.pixels[p][brightnessIndex] = brightnessInt
	return bl
}

func initPixels(brightness int) [8][4]int {
	var pixels [8][4]int
	for p, _ := range pixels {

		pixels[p][redIndex] = 0
		pixels[p][greenIndex] = 0
		pixels[p][blueIndex] = 0
		pixels[p][brightnessIndex] = brightness

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
func NewBlinkt(brightness ...float64) Blinkt {

	//brightness is optional so set the default
	brightnessInt := defaultBrightnessInt

	//over-ride the default if the user has supplied a brightness value
	if len(brightness) > 0 {
		brightnessInt = convertBrightnessToInt(brightness[0])
	}
	return Blinkt{
		pixels: initPixels(brightnessInt),
	}
}

// Blinkt use the NewBlinkt function to initialize the pixels property.
type Blinkt struct {
	pixels [8][4]int
}

func init() {

}
