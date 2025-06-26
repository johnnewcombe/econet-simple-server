package comms

import (
	"context"
	"github.com/johnnewcombe/econet-simple-server/logger"
	"go.bug.st/serial"
	"io"
	"sync"
	"time"
)

// see https://pkg.go.dev/go.bug.st/serial.v1#Mode

type SerialClient struct{}

var (
	port serial.Port
)

func (c *SerialClient) Open(portName string) error {

	var (
		err  error
		mode serial.Mode
	)

	if err = c.Close(); err != nil {
		return err
	}

	// piconet properties
	mode = serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	// Note that he configuration can be changed at any time with the SetMode function:
	if port, err = serial.Open(portName, &mode); err != nil {
		port = nil
		return err
	}

	//port.SetDTR(true)
	//port.SetRTS(true)
	m, err := port.GetModemStatusBits()
	print(m)

	time.Sleep(500 * time.Millisecond)

	return nil
}

func (c *SerialClient) Close() error {

	if port != nil {
		if err := port.Close(); err != nil {
			return err
		}
		port = nil
	}
	return nil
}

func (c *SerialClient) Write(byt []byte) error {

	var err error

	logger.LogDebug.Printf("TX: %s", string(byt))

	if port != nil {
		if _, err = port.Write(byt); err != nil {
			return err
		}
		//print(string(byt))
	}
	return nil

}

func (c *SerialClient) Read(ctx context.Context, wg *sync.WaitGroup, ch chan byte) {

	var (
		sInputByte string
	)

	defer wg.Done()

	for {

		// process any requested cancellation by checking the Done channel of the context
		// note that the resdByte function blocks and so when in that state cancellation
		// wont happen unless there is data coming in or the port is closed
		select {
		case <-ctx.Done():
			// ctx is telling us to stop
			logger.LogDebug.Println("SerialClient.Read() goroutine cancelled.")
			return

		default:
		}

		if port != nil {
			ok, inputByte := c.readByte()

			//logger.LogDebug.Printf("Data Received: %v, Byte: %d\r\n", ok, inputByte)
			if ok {
				// send via a channel

				// Tidy up the logging output
				if inputByte >= 0x20 {
					sInputByte = string(inputByte)
				} else {
					sInputByte = "."
				}
				logger.LogDebug.Printf("RX: %s (%02x)", sInputByte, inputByte)

				// send byte out to the channel, this is blocking until collected
				ch <- inputByte

			}
		} else {
			// no need to rush as the connection isn't open
			time.Sleep(2 * time.Millisecond)
		}
	}
}

func (c *SerialClient) readByte() (bool, byte) {

	result := true
	inputByte := make([]byte, 1, 1)

	// get a byte
	count, err := port.Read(inputByte)

	if err != nil {
		if err != io.EOF {
			return false, 0
		} else {
			if err = c.Close(); err != nil {
				return false, 0
			}
		}
		result = false
	}
	if count == 0 {
		print(count)
	}

	if len(inputByte) > 0 {
		return result, inputByte[0]
	} else {
		return false, 0
	}
}

func (c *SerialClient) GetPortsList() ([]string, error) {
	return serial.GetPortsList()
}
