package main

import (
 . "github.com/alexellis/rpi"
 "fmt"
)

const DAT int = 23
const CLK int = 24

func show(pixels [8][4]int) {
    eof(32)
    for i,_ := range pixels {
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
    eof(36)
}

func initPixels(brightness int) [8][4]int {
    var pixels [8][4] int
    for i,_ := range pixels {
       pixels[i][0] = 0
       pixels[i][1] = 0
       pixels[i][2] = 0
       pixels[i][3] = brightness
    }
    return pixels
}

func eof(pulses int) {
     DigitalWrite(GpioToPin(DAT), 0)
     for i := 0 ; i < pulses ; i++ {
       DigitalWrite(GpioToPin(CLK), 1)
       DigitalWrite(GpioToPin(CLK), 0)
     }
}

func setup() {
  PinMode(GpioToPin(DAT), OUTPUT)
  PinMode(GpioToPin(CLK), OUTPUT)
}

func writeByte(val int) {
  for i := 0; i < 8; i++ {
	// 0b10000000 = 128
     DigitalWrite(GpioToPin(DAT), val & 128)
     DigitalWrite(GpioToPin(CLK), 1)
     val = val << 1
     DigitalWrite(GpioToPin(CLK), 0)
  }
}

func main() {
    WiringPiSetup()

    pixels := initPixels(50)    
    setup()
    fmt.Println(pixels)
    Delay(100)

    pixels[0][1] = 128 & 255 // 0xff = 255
    pixels[1][1] = 128 & 255 // 0xff = 255
    pixels[5][1] = 128 & 255 // 0xff = 255
    pixels[6][1] = 128 & 255 // 0xff = 255

    show(pixels)    
    Delay(1000)
    pixels[0][1] = 0 & 255 // 0xff = 255
    pixels[1][1] = 0 & 255 // 0xff = 255
    pixels[5][1] = 0 & 255 // 0xff = 255
    pixels[6][1] = 0 & 255 // 0xff = 255

    show(pixels)    

}

