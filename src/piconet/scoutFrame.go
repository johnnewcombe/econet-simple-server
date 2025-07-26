package piconet

import "encoding/base64"

type ScoutFrame struct {
	Header
	ControlByte byte
	Port        byte
}

func newScoutFrame(base64EncodedData string) (ScoutFrame, error) {

	var (
		decodedFrame []byte
		err          error
	)

	if decodedFrame, err = base64.StdEncoding.DecodeString(base64EncodedData); err != nil {
		return ScoutFrame{}, err
	}
	if decodedFrame, err = Base64Decode(base64EncodedData); err != nil {
		return ScoutFrame{}, err
	}

	var scout = ScoutFrame{}
	scout.DstStn = decodedFrame[0]
	scout.DstNet = decodedFrame[1]
	scout.SrcStn = decodedFrame[2]
	scout.SrcNet = decodedFrame[3]
	scout.ControlByte = decodedFrame[4]
	scout.Port = decodedFrame[5]

	return scout, nil
}
