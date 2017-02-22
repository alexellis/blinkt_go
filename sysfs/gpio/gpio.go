package gpio

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

const OUTPUT = 1

type gpioPin struct {
	valueFd     *os.File
	directionFd *os.File
}

var gpioPins map[string]gpioPin

func Setup() {
	gpioPins = make(map[string]gpioPin)
}

func Cleanup() {
	for k, v := range gpioPins {
		val, _ := strconv.Atoi(k)
		fmt.Println("Cleaning up " + k)
		v.directionFd.Close()
		v.valueFd.Close()

		unexport(val)
	}
}

func export(pin int) {
	path := "/sys/class/gpio/export"
	bytesToWrite := []byte(strconv.Itoa(pin))
	writeErr := ioutil.WriteFile(path, bytesToWrite, 0644)
	if writeErr != nil {
		log.Panic(writeErr)
	}
}

func unexport(pin int) {
	path := "/sys/class/gpio/unexport"
	bytesToWrite := []byte(strconv.Itoa(pin))
	writeErr := ioutil.WriteFile(path, bytesToWrite, 0644)
	if writeErr != nil {
		log.Panic(writeErr)
	}
}

func pinExported(pin int) bool {
	pinPath := fmt.Sprintf("/sys/class/gpio/gpio%d", pin)
	if file, err := os.Stat(pinPath); err == nil && len(file.Name()) > 0 {
		return true
	}
	return false
}

func PinMode(pin int, val int) {
	pinName := strconv.Itoa(pin)

	exported := pinExported(pin)
	if val == OUTPUT {
		if exported == false {
			export(pin)
		}
	} else {
		if exported == true {
			unexport(pin)
		}
	}

	_, exists := gpioPins[pinName]
	if exists == false {
		pinPath := fmt.Sprintf("/sys/class/gpio/gpio%d", pin)
		valueFd, openErr := os.OpenFile(pinPath+"/value", os.O_WRONLY, 0640)
		if openErr != nil {
			log.Panic(openErr, pinPath)
		}
		directionFd, openErr := os.OpenFile(pinPath+"/direction", os.O_WRONLY, 0640)
		if openErr != nil {
			log.Panic(openErr, pinPath)
		}
		gpioPins[pinName] = gpioPin{
			valueFd:     valueFd,
			directionFd: directionFd,
		}
		if val == OUTPUT {
			pinDigitalWrite(pin, "out", "direction")
		}
	}
}

func GpioToPin(pin int) int {
	return pin
}

func DigitalWrite(pin int, val int) {
	pinDigitalWrite(pin, strconv.Itoa(val), "value")
}

func pinDigitalWrite(pin int, val string, mode string) {
	pinName := strconv.Itoa(pin)
	var err error
	if mode == "direction" {
		_, err = gpioPins[pinName].directionFd.Write([]byte(val))
	} else {
		_, err = gpioPins[pinName].valueFd.Write([]byte(val))
	}

	if err != nil {
		log.Panic(err, fmt.Sprintf("Pin: %s Mode: %s Value: %s ", pinName, val, mode))
	}
}

