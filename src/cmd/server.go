package cmd

import (
	"context"
	"fmt"
	"github.com/johnnewcombe/econet-simple-server/src/admin"
	comms2 "github.com/johnnewcombe/econet-simple-server/src/comms"
	piconet "github.com/johnnewcombe/econet-simple-server/src/piconet"
	"github.com/johnnewcombe/econet-simple-server/src/server"
	"github.com/johnnewcombe/econet-simple-server/src/utils"
	"github.com/spf13/cobra"
	"log"
	"log/slog"
	"os"
	"sync"
)

var fileserver = &cobra.Command{
	Use:   "server",
	Short: "Starts the Econet file server.",
	Long: `
Starts the Econet file server.
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		const (
			kPasswordFile = "PASSWORD"
		)
		var (
			err            error
			commsClient    comms2.CommunicationClient
			ctxCommsClient context.Context
			wgComms        sync.WaitGroup
			cancelRead     context.CancelFunc
			portName       string
			debug          bool
			rootFolder     string
			userData       string
			rxChannel      chan byte
		)
		// TODO put the debug in a more generic place e.g. Root Cmd
		if debug, err = cmd.Flags().GetBool("debug"); err != nil {
			return err
		}

		if portName, err = cmd.Flags().GetString("port"); err != nil {
			return err
		}
		if rootFolder, err = cmd.Flags().GetString("root-folder"); err != nil {
			return err
		}

		if debug {
			slog.SetLogLoggerLevel(slog.LevelDebug)
			log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
		} else {
			slog.SetLogLoggerLevel(slog.LevelInfo)
			log.SetFlags(log.Ldate | log.Lmicroseconds)
		}

		var pwFile = rootFolder + "/" + kPasswordFile

		// create a serial client
		commsClient = &comms2.SerialClient{}

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

		//sort root folder
		slog.Info("Opening root folder.", "root-folder", rootFolder)

		if err = utils.CreateDirectoryIfNotExists(rootFolder); err != nil {
			return err
		}

		// check for password file
		slog.Info("Checking for password file.", "password-file", pwFile)
		if !utils.Exists(pwFile) {

			// create new file
			user := admin.User{
				Username:   "SYST",
				Password:   "SYST",
				FreeSpace:  1024e3,
				BootOption: 0b00000000,
				Privilege:  0b11000000,
			}

			users := admin.Users{
				Users: []admin.User{user},
			}

			s := users.ToString()
			if err = utils.WriteString(pwFile, s); err != nil {
				return err
			}
		}

		slog.Info("Loading password file.", "password-file", pwFile)
		userData, err = utils.ReadString(pwFile)
		if err != nil {
			return err
		}

		// load the users
		var users = admin.Users{}
		if err = users.Load(userData); err != nil {
			return err
		}

		// open the port to the piconet device
		slog.Info("Opening serial port.", "port", portName)
		if err = openConnection(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		//var ports, _ = commsClient.GetPortsList()
		//print(ports)
		piconet.SetStationID(commsClient, 254)
		piconet.SetMode(commsClient, "LISTEN")
		//piconet.SetMode(commsClient, "MONITOR")
		piconet.Status(commsClient)

		// start the server
		slog.Info("Listening.", "port-name", portName)

		server.Listener(commsClient, rxChannel)
		slog.Info("No longer listening.", "port-name", portName)

		// server shutdown
		piconet.SetMode(commsClient, "STOP")

		slog.Info("Closing port.", "port-name", portName)

		if err = closeConnection(); err != nil {
			return err
		}

		slog.Info("Server shutdown.")

		return nil
	},
}
