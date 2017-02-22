blinkt_go
=========

> This library is work-in-progress and makes use of WiringPi and the Golang `rpi` library.

## Instructions:

### Install Go

If you don't have Go on the Pi, download it from: https://golang.org/dl/ and pick the armv6l edition.

```
sudo tar -xvf go1.7.4.linux-armv6l.tar.gz -C /usr/local/
export GOPATH=$HOME/go
```

### Install WiringPi
```
# sudo apt-get install -qy wiringpi
```

### Install and build blinkt! library

```
# export GOPATH=$HOME/go/
# mkdir -p $GOPATH/src/github.com/alexellis/
# cd $GOPATH/src/github.com/alexellis/

# git clone https://github.com/alexellis/blinkt_go && cd blinkt_go

# go get
# go build
# sudo ./blinkt_go
```

## sysfs implementation and Docker Swarm

Docker Swarm cannot run the main version of this library because WiringPi needs elevated privileges.

To use this library with Docker Swarm please use the version of Blinkt in the sysfs package.

> For more information on sysfs see: http://elinux.org/RPi_GPIO_Code_Samples#sysfs


## Related:

* [Blinkt Golang examples programs](https://github.com/alexellis/blinkt_go_examples)
* [Golang rpi library](https://github.com/alexellis/rpi/)
