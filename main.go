package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	// "github.com/huin/goserial"
	"go.bug.st/serial.v1"
)

// findArduino looks for the file that represents the Arduino
// serial connection. Returns the fully qualified path to the
// device if we are able to find a likely candidate for an
// Arduino, otherwise an empty string if unable to find
// something that 'looks' like an Arduino device.
func findArduino() string {
	contents, _ := ioutil.ReadDir("/dev")

	// Look for what is mostly likely the Arduino device
	for _, f := range contents {
		if strings.Contains(f.Name(), "tty.usbserial") ||
			strings.Contains(f.Name(), "ttyUSB") {
			return "/dev/" + f.Name()
		}
	}

	// Have not been able to find a USB device that 'looks'
	// like an Arduino.
	return ""
}

// func main() {
// 	// Find the device that represents the arduino serial
// 	// connection.
// 	c := &goserial.Config{Name: findArduino(), Baud: 9600}
// 	_, err := goserial.OpenPort(c)
// 	if err != nil {
// 		fmt.Print(err)
// 	}
// }

func main() {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}

	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open("/dev/cu.usbserial-14340", mode)
	if err != nil {
		log.Fatal(err)
	}
	n, err := port.Write([]byte("echo"))
	// n, err := port.Write([]byte{0x65, 0x63, 0x68, 0x6f})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)

	buff := make([]byte, 100)
	for {
		n, err := port.Read(buff)
		if err != nil {
			log.Fatal(err)
			break
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}
		fmt.Printf("%v", hex.Dump(buff[:n]))
	}
}
