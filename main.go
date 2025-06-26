package main

import (
	"context"
	"fmt"
	"github.com/johnnewcombe/econet-simple-server/comms"
	"github.com/johnnewcombe/econet-simple-server/logger"
	"github.com/johnnewcombe/econet-simple-server/piconet"
	"github.com/johnnewcombe/econet-simple-server/server"
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
		rxChannel      chan byte
	)

	portName = "/dev/tty.usbmodem14301" // TODO Get from config or auto detect
	//portName = "/dev/tty.usbserial-1440"

	// create a serial client
	commsClient = &comms.SerialClient{}

	// create the channel that will receive the data
	rxChannel = make(chan byte)

	// create a wait group and make sure we wait for all goroutines to end before exiting
	wgComms = sync.WaitGroup{}

	// define the Open function
	openConnection := func() error {

		if err = commsClient.Open(portName); err != nil {

			return err
		}

		// move cursor down a line, makes for better output
		fmt.Println()

		logger.LogInfo.Printf("Opening Port: %s", portName)

		// about to start a go routine so add 1 to the waitgroup
		wgComms.Add(1)
		ctxCommsClient, cancelRead = context.WithCancel(context.Background())

		logger.LogInfo.Printf("Listening on port: %s", portName)

		// start the read go routine passing in the rx channel on which to return data
		// the data is collected by the listener function
		go commsClient.Read(ctxCommsClient, &wgComms, rxChannel)

		// all done opening the port
		return nil
	}

	closeConnection := func() {
		// close the client and stop the read goroutine.
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

	openConnection()
	//var ports, _ = commsClient.GetPortsList()
	//print(ports)

	// initialisation
	piconet.SetStationID(commsClient, 254)
	piconet.SetMode(commsClient, "LISTEN")
	piconet.GetStatus(commsClient)

	// start the server
	server.Listener(rxChannel)

	logger.LogInfo.Printf("No longer listening on port: %s", portName)

	// server shutdown
	piconet.SetMode(commsClient, "STOP")

	logger.LogInfo.Printf("Closing port: %s", portName)
	closeConnection()
	logger.LogInfo.Println("Server shutdowm.")

}
