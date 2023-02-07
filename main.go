package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"go.bug.st/serial.v1"
)

var (
    port  *string
)

func init() {
    port = flag.String("port", "", "serial port")
}


func NewStreamer(rwc io.ReadWriteCloser) (<-chan []byte, chan<- []byte) {
	inCh := make(chan []byte)
	// outCh := make(chan []byte)

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := rwc.Read(buf)

			// This error should be handled somehow
			_ = err

			// This buffer is now not safe to reuse again because the other side of the channel has it.
			// We should make a copy, but that adds an extra allocation.
			// outCh <- buf[:n]
			fmt.Printf("%v", hex.Dump(buf[:n]))
		}
	}()

	go func() {
		for {
			buf := <-inCh
			n, err := rwc.Write(buf)

			// We should return the result to the caller somehow.
			// Or at least handle an error
			_ = n
			_ = err
		}
	}()

	return nil, inCh
}

func main() {
	flag.Parse()
	if *port == "" {
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
		return
	} else {
		fmt.Println("selected port:", *port)
	}

	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(*port, mode)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	// inCh, outCh := NewStreamer(port)
	_, outCh := NewStreamer(port)

	// go func() {
	// 	for {
	// 		fmt.Printf("%v", hex.Dump(<-inCh))
	// 	}
	// }()

	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		switch text {
		case "find1":
			outCh <- []byte{0x05, 0xFA, 0x06, 0xF9, 0x07, 0xF8, 0x00, 0xFF}
			fmt.Println("Send find_1 controller sequence")
		case "find2":
			// port.SetRTS(true)
			// port.SetDTR(true)
			outCh <- []byte{0xfa, 0x06, 0xf9, 0x07, 0xf8, 0x00, 0xff}
			fmt.Println("Send find_2 controller sequence")
		case "find3":
			outCh <- []byte{0x06, 0xF9, 0x07, 0xF8, 0x00, 0xFF}
			fmt.Println("Send find_3 controller sequence")
		case "find4":
			outCh <- []byte{0xF9, 0x07, 0xF8, 0x00, 0xFF}
			fmt.Println("Send find_4 controller sequence")
		case "find5":
			outCh <- []byte{0x07, 0xF8, 0x00, 0xFF}
			fmt.Println("Send find_5 controller sequence")
		case "find6":
			outCh <- []byte{0xF8, 0x00, 0xFF}
			fmt.Println("Send find_6 controller sequence")
		case "find7":
			outCh <- []byte{0x00, 0xFF}
			fmt.Println("Send find_7 controller sequence")
		case "find8":
			outCh <- []byte{0xFF}
			fmt.Println("Send find_8 controller sequence")
		case "init":
			outCh <- []byte{0x07, 0xF8, 0x00, 0xFF}
			fmt.Println("Send init controller sequence")
		case "start":
			outCh <- []byte{
				0x01, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
				0x00, 0xfe, 0x02, 0x0a, 0xbc, 0xfd, 
			}
			fmt.Println("Send start controller sequence")
		case "stop":
			outCh <- []byte{
				0x01, 0xFF,
			}
			fmt.Println("Send stop controller sequence")
		case "a":
			outCh <- []byte{
				0x01, 0x41, 0x72, 0x64, 0x75, 0x69, 0x6e, 0x6f, 0x20, 0x4b, 0x6e, 0x69, 0x74, 0x74, 0x20, 0x4d, 
				0x61, 0x63, 0x68, 0x69, 0x6e, 0x65, 0x20, 0x61, 0x64, 0x61, 0x70, 0x74, 0x65, 0x72, 0x0a, 0x47, 
				0x72, 0x61, 0x79, 0x59, 0x6f, 0x67, 0x69, 0x20, 0x28, 0x63, 0x29, 0x20, 0x32, 0x30, 0x31, 0x39, 
				0x0a, 0x00, 0xf7, 0x0d,
			}
			fmt.Println("Send answer controller sequence")
		default:
			outCh <- []byte(text) // just send
			fmt.Printf("Just send [%s]\n", text)
		}

	}
}
