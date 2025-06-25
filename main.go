package main

import (
	"context"
	"fmt"
	"github.com/johnnewcombe/econet-simple-server/comms"
	"sync"
)

func main() {

	var (
		err            error
		commsClient    comms.CommunicationClient
		ctxCommsClient context.Context
		wgComms        sync.WaitGroup
		cancelRead     context.CancelFunc
		portName       string
	)

	portName = "/dev/tty.usbmodem14301"
	commsClient = &comms.SerialClient{}

	// create a wait group and make sure we wait for all goroutines to end before exiting
	wgComms = sync.WaitGroup{}

	// this needs to be here in case the initial open fails, and the user selects another
	//ctxCommsClient, cancelRead = context.WithCancel(context.Background())
	// define the Open function
	openFunc := func() error {

		if err = commsClient.Open(portName); err != nil {

			return err
		}

		wgComms.Add(1)
		ctxCommsClient, cancelRead = context.WithCancel(context.Background())
		go commsClient.Read(ctxCommsClient, &wgComms, func(ok bool, b byte) {
			if ok {
				fmt.Println(b)
			} else {
				fmt.Println("Offline")
			}
		})
		fmt.Printf("Offline %s\r\n", portName)
		return nil
	}

	closeFunc := func() {
		// close the previous client and stop the read goroutine.
		// The commsClient.Read() goroutine blocks on serial/net. Closing the
		// connection/port will cause a read error and allow the go routine to continue
		// monitoring for ctx cancel
		commsClient.Close()

		// (it raises and error instead) and it is now looping until cancelled,
		// so lets cancel it
		if cancelRead != nil {
			cancelRead()
		}

		// wait for all goroutines to stop
		wgComms.Wait()
	}

	exitFunc := func() {

		// order is important
		commsClient.Close()
		if cancelRead != nil {
			cancelRead()
		}
		wgComms.Wait()
	}

	openFunc()

	commsClient.Write([]byte("SET_STATION 121\r"))
	//commsClient.Write([]byte("SET_MODE MONITOR\r"))

	closeFunc()
	exitFunc()
}

// These are all piconet commands not Econet ones
// "SET_MODE STOP\r"
// "SET_MODE MONITOR\r"
// "SET_MODE LISTEN\r"
// "SET_STATION 121\r"
