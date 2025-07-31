package cobra

import (
	"context"
	"fmt"
	"github.com/johnnewcombe/econet-simple-server/src/comms"
	"github.com/johnnewcombe/econet-simple-server/src/econet"
	"github.com/johnnewcombe/econet-simple-server/src/piconet"
	"github.com/johnnewcombe/econet-simple-server/src/server"
	"github.com/johnnewcombe/econet-simple-server/src/utils"
	"github.com/spf13/cobra"
	"log"
	"log/slog"
	"os"
	"sync"
)

var fileserver = &cobra.Command{
	TraverseChildren: true,
	Use:              "server",
	Short:            "Starts the Econet file server.",
	Long: `
Starts the Econet file server.
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			err            error
			commsClient    comms.CommunicationClient
			ctxCommsClient context.Context
			wgComms        sync.WaitGroup
			cancelRead     context.CancelFunc
			portName       string
			debug          bool
			rootFolder     string
			rxChannel      chan byte
			users          econet.Passwords
		)
		// TODO put the debug in a more generic place e.g. Root Event
		// get data passed in via flags
		if debug, err = cmd.Flags().GetBool("debug"); err != nil {
			return err
		}
		if portName, err = cmd.Flags().GetString("port"); err != nil {
			return err
		}
		if rootFolder, err = cmd.Flags().GetString("root-folder"); err != nil {
			return err
		}

		// configure logging
		if debug {
			slog.SetLogLoggerLevel(slog.LevelDebug)
			log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
		} else {
			slog.SetLogLoggerLevel(slog.LevelInfo)
			log.SetFlags(log.Ldate | log.Lmicroseconds)
		}

		//sort root folder
		slog.Info("Checking the root folder.", "root-folder", rootFolder)

		// set globals
		econet.LocalRootDiectory = rootFolder + "/"
		econet.LocalDisk0 = econet.LocalRootDiectory + econet.Disk0
		econet.LocalDisk1 = econet.LocalRootDiectory + econet.Disk1
		econet.LocalDisk2 = econet.LocalRootDiectory + econet.Disk2
		econet.LocalDisk3 = econet.LocalRootDiectory + econet.Disk3

		// cteate directories if needed
		if err = utils.CreateDirectoryIfNotExists(econet.LocalRootDiectory); err != nil {
			return err
		}
		if err = utils.CreateDirectoryIfNotExists(econet.LocalRootDiectory + econet.Disk0); err != nil {
			return err
		}
		if err = utils.CreateDirectoryIfNotExists(econet.LocalRootDiectory + econet.Disk1); err != nil {
			return err
		}
		if err = utils.CreateDirectoryIfNotExists(econet.LocalRootDiectory + econet.Disk2); err != nil {
			return err
		}
		if err = utils.CreateDirectoryIfNotExists(econet.LocalRootDiectory + econet.Disk3); err != nil {
			return err
		}

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

		// check for password file
		var pwFile = econet.LocalRootDiectory + econet.PasswordFile

		slog.Info("Checking for password file.", "password-file", pwFile)

		if !utils.Exists(pwFile) {

			slog.Info("Creating new password file.", "password-file", pwFile)
			// create new file
			user := econet.User{
				Username:   "SYST",
				Password:   "SYST",
				FreeSpace:  1024e3,
				BootOption: 0b00000000,
				Privilege:  0b11000000,
			}

			// add the user to the userData
			userData := econet.Passwords{
				Items: []econet.User{user},
			}

			// write the userData to disk
			s := userData.ToString()
			if err = utils.WriteString(pwFile, s); err != nil {
				return err
			}
		}

		// load the users data as a sanity check
		if users, err = econet.NewUsers(pwFile); err != nil {
			return err
		}

		// store the users data in the global variable
		econet.Userdata = users

		slog.Info("Password file valid.", "password-file", pwFile, "user-count", len(users.Items))

		// open the port to the piconet device
		slog.Info("Opening serial port.", "port", portName)
		if err = openConnection(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		// TODO: could this be used for an 'auto' port select
		//var ports, _ = commsClient.GetPortsList()
		//print(ports)
		piconet.SetStationID(commsClient, 254)
		piconet.SetMode(commsClient, "LISTEN")
		piconet.Status(commsClient)

		// start the server
		server.Listener(commsClient, rxChannel)
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
