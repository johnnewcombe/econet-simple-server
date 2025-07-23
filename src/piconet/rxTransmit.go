package piconet

import (
	"encoding/base64"
	"fmt"
	"strings"
)

type Header struct {
	DstStn      byte
	DstNet      byte
	SrcStn      byte
	SrcNet      byte
	ControlByte byte
	Port        byte
}

type ScoutFrame struct {
	Header
}
type DataFrame struct {
	Header
	Data []byte
}

type RxTransmit struct {
	ScoutFrame ScoutFrame
	DataFrame  DataFrame
}

func NewRxTransmit(command Cmd) (RxTransmit, error) {

	var (
		err   error
		rxt   RxTransmit
		scout ScoutFrame
		data  DataFrame
	)

	if scout, err = newScout(command.Args[0]); err != nil {
		return RxTransmit{}, err
	}

	if data, err = newData(command.Args[1]); err != nil {
		return RxTransmit{}, err
	}

	rxt = RxTransmit{
		ScoutFrame: scout,
		DataFrame:  data,
	}

	return rxt, nil
}

func (rxt *RxTransmit) String() string {

	var sb = strings.Builder{}
	sb.WriteString(fmt.Sprintf("scout-dst-stn=%02X, scout-dst-net=%02X, scout-src-stn=%02X, scout-scr-net=%02X, scout-ctrl-byte=%02X, scout-port=%02X, scout-port-desc=%s",
		rxt.ScoutFrame.DstStn, rxt.ScoutFrame.DstNet, rxt.ScoutFrame.SrcStn, rxt.ScoutFrame.SrcNet, rxt.ScoutFrame.ControlByte, rxt.ScoutFrame.Port, PortMap[rxt.ScoutFrame.Port]))
	sb.WriteString(", ")
	sb.WriteString(fmt.Sprintf("data-dst-stn=%02X, data-dst-net=%02X, data-src-stn=%02X, data-scr-net=%02X, data-ctrl-byte=%02X, data-port=%02X, data-port-desc=%s",
		rxt.DataFrame.DstStn, rxt.DataFrame.DstNet, rxt.DataFrame.SrcStn, rxt.DataFrame.SrcNet, rxt.DataFrame.ControlByte, rxt.DataFrame.Port, PortMap[rxt.DataFrame.Port]))

	if len(rxt.DataFrame.Data) > 0 {
		sb.WriteString(fmt.Sprintf(", data=[% 02X]", rxt.DataFrame.Data))
	}

	return sb.String()
}

func newScout(base64EncodedData string) (ScoutFrame, error) {

	var (
		decodedFrame []byte
		err          error
	)

	if decodedFrame, err = base64.StdEncoding.DecodeString(base64EncodedData); err != nil {
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
func newData(base64EncodedData string) (DataFrame, error) {

	var (
		decodedFrame []byte
		err          error
	)

	if decodedFrame, err = base64.StdEncoding.DecodeString(base64EncodedData); err != nil {
		return DataFrame{}, err
	}

	var data = DataFrame{}
	data.DstStn = decodedFrame[0]
	data.DstNet = decodedFrame[1]
	data.SrcStn = decodedFrame[2]
	data.SrcNet = decodedFrame[3]
	data.ControlByte = decodedFrame[4]
	data.Port = decodedFrame[5]

	// Add any data (Scouts don't have data)
	if len(decodedFrame) > 6 {
		data.Data = decodedFrame[6:]
	}
	return data, nil
}
