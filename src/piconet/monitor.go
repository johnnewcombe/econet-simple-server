package piconet

import "encoding/base64"

type Monitor struct {
	Frame []byte
}

func NewMonitor(eventArgs []string) (Monitor, error) {

	var (
		decodedFrame []byte
		err          error
	)

	if decodedFrame, err = base64.StdEncoding.DecodeString(eventArgs[0]); err != nil {
		return Monitor{}, err
	}

	return Monitor{
		Frame: decodedFrame,
	}, nil

}
