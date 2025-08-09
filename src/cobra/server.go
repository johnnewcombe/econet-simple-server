package cobra

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/johnnewcombe/econet-simple-server/src/econet"
	"github.com/johnnewcombe/econet-simple-server/src/lib"
	"github.com/johnnewcombe/econet-simple-server/src/piconet"
	"github.com/spf13/cobra"
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
			commsClient    piconet.CommunicationClient
			ctxCommsClient context.Context
			wgComms        sync.WaitGroup
			cancelRead     context.CancelFunc
			portName       string
			debug          bool
			rootFolder     string
			rxChannel      chan byte
			users          econet.Users
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
		econet.LocalDisk0 = econet.LocalRootDiectory + econet.Disk0 + "/"
		econet.LocalDisk1 = econet.LocalRootDiectory + econet.Disk1 + "/"
		econet.LocalDisk2 = econet.LocalRootDiectory + econet.Disk2 + "/"
		econet.LocalDisk3 = econet.LocalRootDiectory + econet.Disk3 + "/"

		// cteate directories if needed
		if err = lib.CreateDirectoryIfNotExists(econet.LocalRootDiectory); err != nil {
			return err
		}
		if err = lib.CreateDirectoryIfNotExists(econet.LocalRootDiectory + econet.Disk0); err != nil {
			return err
		}
		if err = lib.CreateDirectoryIfNotExists(econet.LocalRootDiectory + econet.Disk1); err != nil {
			return err
		}
		if err = lib.CreateDirectoryIfNotExists(econet.LocalRootDiectory + econet.Disk2); err != nil {
			return err
		}
		if err = lib.CreateDirectoryIfNotExists(econet.LocalRootDiectory + econet.Disk3); err != nil {
			return err
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

		// check for a password file
		var pwFile = econet.LocalRootDiectory + econet.PasswordFile

		slog.Info("Checking for password file.", "password-file", pwFile)

		// load the user data as a sanity check
		if users, err = econet.NewUsers(pwFile); err != nil {
			return err
		}

		// ensure SYST has a home directorey, all home directories are on Disk 0
		if err = lib.CreateDirectoryIfNotExists(econet.LocalDisk0 + econet.DefaultSystemUserName); err != nil {
			return err
		}

		// store the user data in the global variable
		econet.Userdata = users

		slog.Info("Password file valid.", "password-file", pwFile, "user-count", len(users.Items))

		// open the port to the piconet device
		slog.Info("Opening device.", "device-name", portName)
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
		piconet.Listener(commsClient, rxChannel)

		slog.Info("No longer listening.", "device-name", portName)

		// server shutdown
		piconet.SetMode(commsClient, "piconet-cmd=STOP")

		slog.Info("Closing device.", "device-name", portName)

		if err = closeConnection(); err != nil {
			return err
		}

		slog.Info("Server shutdown.")

		return nil
	},
}
