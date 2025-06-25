package main

import (
	"comms"
	"fmt"
	"sync"
)

func main() {

	var (
		wgComms sync.WaitGroup
	)

	// create a wait group and make sure we wait for all goroutines to end before exiting
	wgComms = sync.WaitGroup{}

	// this needs to be here in case the initial open fails, and the user selects another
	//ctxCommsClient, cancelRead = context.WithCancel(context.Background())
	// define the Open function
	openFunc := func() error {

		// create the communications client
		if endpoint.IsSerial() {
			commsClient = &comms.SerialClient{}
		} else {
			commsClient = &comms.NetClient{}
		}

		if err = commsClient.Open(endpoint); err != nil {

			return err
		}

		wgComms.Add(1)
		ctxCommsClient, cancelRead = context.WithCancel(context.Background())
		go commsClient.Read(ctxCommsClient, &wgComms, func(ok bool, b byte) {
			if ok {
				if err = screen.Write(b); err != nil {
					// It is safe to updated Fyne UI components from within a go routine, in fact its probably the only
					// way to do it and is sited in Fyne examples. See https://developer.fyne.io/started/updating
					displayError(err, w)
					return
				}
			} else {
				status.Text = "Offline"
				status.Refresh()
			}
		})
		status.Text = fmt.Sprintf("Online [%s]", endpoint.Name)
		status.Refresh()
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

}
