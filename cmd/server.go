package cmd

import (
	"context"
	"fmt"
	"github.com/johnnewcombe/econet-simple-server/admin"
	"github.com/johnnewcombe/econet-simple-server/comms"
	"github.com/johnnewcombe/econet-simple-server/logger"
	"github.com/johnnewcombe/econet-simple-server/piconet"
	"github.com/johnnewcombe/econet-simple-server/server"
	"github.com/johnnewcombe/econet-simple-server/utils"
	"github.com/spf13/cobra"
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
			commsClient    comms.CommunicationClient
			ctxCommsClient context.Context
			wgComms        sync.WaitGroup
			cancelRead     context.CancelFunc
			portName       string
			rootFolder     string
			userData       string
			rxChannel      chan byte
		)

		if portName, err = cmd.Flags().GetString("port"); err != nil {
			return err
		}
		if rootFolder, err = cmd.Flags().GetString("root-folder"); err != nil {
			return err
		}
		var pwFile = rootFolder + "/" + kPasswordFile

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

		//sort root folder
		logger.LogInfo.Printf("Opening Root Folder: %s", rootFolder)
		if err = utils.CreateDirectoryIfNotExists(rootFolder); err != nil {
			return err
		}

		// check for password file
		logger.LogInfo.Printf("Checking for Password file: %s", pwFile)
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

		logger.LogInfo.Printf("Loading Password file: %s", pwFile)
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
		logger.LogInfo.Printf("Opening Port: %s", portName)
		if err = openConnection(); err != nil {
			logger.LogError.Fatalf("%s (%s)", err, portName)
		}

		//var ports, _ = commsClient.GetPortsList()
		//print(ports)

		// initialisation
		piconet.SetStationID(commsClient, 254)
		piconet.SetMode(commsClient, "LISTEN")
		piconet.GetStatus(commsClient)

		// start the server
		logger.LogInfo.Printf("Listening on port: %s", portName)
		server.Listener(commsClient, rxChannel)
		logger.LogInfo.Printf("No longer listening on port: %s", portName)

		// server shutdown
		piconet.SetMode(commsClient, "STOP")

		logger.LogInfo.Printf("Closing port: %s", portName)

		if err = closeConnection(); err != nil {
			return err
		}

		logger.LogInfo.Println("Server shutdown.")

		return nil
	},
}
