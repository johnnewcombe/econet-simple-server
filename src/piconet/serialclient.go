package piconet

import (
	"context"
	"fmt"
	"go.bug.st/serial"
	"io"
	"log/slog"
	"sync"
	"time"
)

// see https://pkg.go.dev/go.bug.st/serial.v1#Mode
type CommunicationClient interface {
	Open(portName string) error
	Write(byt []byte) error
	Read(ctx context.Context, wg *sync.WaitGroup, c chan byte)
	Close() error
}

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

	// Note that the configuration can be changed at any time with the SetMode function:
	if port, err = serial.Open(portName, &mode); err != nil {
		port = nil
		return err
	}

	//port.SetDTR(true)
	//port.SetRTS(true)
	//m, err := port.GetModemStatusBits()

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

	var (
		err error
	)

	if port != nil {
		// send a byte at a time for logging purposes
		for _, b := range byt {

			if _, err = port.Write([]byte{b}); err != nil {
				return err
			}

			slog.Debug(fmt.Sprintf("tx-ascii=%s, tx-hex=%02X", logTidy(b), b))
		}

	}
	return nil

}

func (c *SerialClient) Read(ctx context.Context, wg *sync.WaitGroup, ch chan byte) {

	defer wg.Done()

	for {

		// process any requested cancellation by checking the Done channel of the context
		// note that the resdByte function blocks and so when in that state cancellation
		// wont happen unless there is data coming in or the port is closed
		select {
		case <-ctx.Done():
			// ctx is telling us to stop
			slog.Debug("SerialClient.Read() goroutine cancelled.")
			return

		default:
		}

		if port != nil {
			ok, b := c.readByte()

			//logger.LogDebug.Printf("DataFrame Received: %v, Byte: %d\r\n", ok, inputByte)
			if ok {

				slog.Debug(fmt.Sprintf("rx-ascii=%s,rx-hex=%02X", logTidy(b), b))

				// send byte out to the channel, this is blocking until collected
				ch <- b

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

func logTidy(b byte) string {
	// tidy up for logging
	if b < 0x20 {
		return "."
	} else {
		return string(b)
	}
}
