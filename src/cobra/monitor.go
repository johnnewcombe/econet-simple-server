package cobra

import (
	"context"
	"fmt"
	"github.com/johnnewcombe/econet-simple-server/src/piconet"
	"github.com/spf13/cobra"
	"log"
	"log/slog"
	"os"
	"sync"
)

var monitor = &cobra.Command{
	Use:   "monitor",
	Short: "Starts the Econet file server in monitor mode.",
	Long: `
Starts the Econet file server in monitor mode.
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			err            error
			commsClient    piconet.CommunicationClient
			ctxCommsClient context.Context
			wgComms        sync.WaitGroup
			cancelRead     context.CancelFunc
			portName       string
			debug          bool
			rxChannel      chan byte
		)

		if debug, err = cmd.Flags().GetBool("debug"); err != nil {
			return err
		}

		if portName, err = cmd.Flags().GetString("port"); err != nil {
			return err
		}

		if debug {
			slog.SetLogLoggerLevel(slog.LevelDebug)
			log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
		} else {
			slog.SetLogLoggerLevel(slog.LevelInfo)
			log.SetFlags(log.Ldate | log.Lmicroseconds)
		}

		// create a serial client
		commsClient = &piconet.SerialClient{}

		// create the channel that will receive the data
		rxChannel = make(chan byte)

		// create a wait group and make sure we wait for all goroutines to end before exiting
		wgComms = sync.WaitGroup{}

		// define the Open function
		openConnection := func() error {

			if err = commsClient.Open(portName); err != nil {

				return err
			}

			// about to start a go routine so add 1 to the waitGroup
			wgComms.Add(1)
			ctxCommsClient, cancelRead = context.WithCancel(context.Background())

			// start the read go routine passing in the rx channel on which to return data
			// the data is collected by the listener function
			go commsClient.Read(ctxCommsClient, &wgComms, rxChannel)

			// all done opening the port
			return nil
		}

		closeConnection := func() error {
			// close the client and stop the read goroutine.
			// The commsClient.Read() goroutine blocks on serial/net. Closing the
			// connection/port will cause a read error and allow the go routine to continue
			// monitoring for ctx cancel
			if err = commsClient.Close(); err != nil {
				return err
			}

			// (it raises and error instead) and it is now looping until cancelled,
			// so lets cancel it
			if cancelRead != nil {
				cancelRead()
			}

			// wait for all goroutines to stop
			wgComms.Wait()

			return nil
		}

		// move cursor down a line, makes for better output
		fmt.Println()

		// open the port to the piconet device
		slog.Info("Opening serial port.", "port", portName)
		if err = openConnection(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		//var ports, _ = commsClient.GetPortsList()
		//print(ports)
		piconet.SetStationID(commsClient, 254)
		piconet.SetMode(commsClient, "MONITOR")
		//piconet.SetMode(commsClient, "MONITOR")
		piconet.Status(commsClient)

		// start the server
		piconet.Listener(commsClient, rxChannel)
		slog.Info("No longer listening.", "port-name", portName)

		// server shutdown
		piconet.SetMode(commsClient, "piconet-cmd=STOP")

		slog.Info("Closing port.", "port-name", portName)

		if err = closeConnection(); err != nil {
			return err
		}

		slog.Info("Server shutdown.")

		return nil
	},
}
