package piconet

import (
	"encoding/base64"
	"fmt"
	"github.com/johnnewcombe/econet-simple-server/src/econet"
	"strings"
)

type RxTransmit struct {
	ScoutFrame *econet.ScoutFrame
	DataFrame  *econet.DataFrame
}

func NewRxTransmit(eventArgs []string) (*RxTransmit, error) {

	var (
		err   error
		rxt   RxTransmit
		scout *econet.ScoutFrame
		data  *econet.DataFrame
	)

	if scout, err = newScoutFrame(eventArgs[0]); err != nil {
		return &RxTransmit{}, err
	}

	if data, err = newDataFrame(eventArgs[1]); err != nil {
		return &RxTransmit{}, err
	}

	rxt = RxTransmit{
		ScoutFrame: scout,
		DataFrame:  data,
	}

	return &rxt, nil
}

func (rxt *RxTransmit) String() string {

	var sb = strings.Builder{}
	sb.WriteString(fmt.Sprintf("scout-dst-stn=%02X, scout-dst-net=%02X, scout-src-stn=%02X, scout-scr-net=%02X, scout-ctrl-byte=%02X, scout-port=%02X, scout-port-desc=%s",
		rxt.ScoutFrame.DstStn, rxt.ScoutFrame.DstNet, rxt.ScoutFrame.SrcStn, rxt.ScoutFrame.SrcNet, rxt.ScoutFrame.ControlByte, rxt.ScoutFrame.Port, PortMap[rxt.ScoutFrame.Port]))
	sb.WriteString(", ")
	sb.WriteString(fmt.Sprintf("data-dst-stn=%02X, data-dst-net=%02X, data-src-stn=%02X, data-scr-net=%02X, data-ctrl-byte=%02X",
		rxt.DataFrame.DstStn, rxt.DataFrame.DstNet, rxt.DataFrame.SrcStn, rxt.DataFrame.SrcNet))

	if len(rxt.DataFrame.Data) > 0 {
		sb.WriteString(fmt.Sprintf(", data-bytes=[% 02X]", rxt.DataFrame.Data))
	}

	return sb.String()
}

func (rxt *RxTransmit) Command() string {

	if len(rxt.DataFrame.Data) > 4 {
		return strings.TrimRight(string(rxt.DataFrame.Data[4:]), "\r")
	}
	return ""
}

func newScoutFrame(base64EncodedData string) (*econet.ScoutFrame, error) {

	var (
		decodedFrame []byte
		err          error
	)

	if decodedFrame, err = base64.StdEncoding.DecodeString(base64EncodedData); err != nil {
		return &econet.ScoutFrame{}, err
	}
	if decodedFrame, err = Base64Decode(base64EncodedData); err != nil {
		return &econet.ScoutFrame{}, err
	}

	var scout = econet.ScoutFrame{}
	scout.DstStn = decodedFrame[0]
	scout.DstNet = decodedFrame[1]
	scout.SrcStn = decodedFrame[2]
	scout.SrcNet = decodedFrame[3]
	scout.ControlByte = decodedFrame[4]
	scout.Port = decodedFrame[5]

	return &scout, nil
}

func newDataFrame(base64EncodedData string) (*econet.DataFrame, error) {

	var (
		decodedFrame []byte
		err          error
	)

	if decodedFrame, err = base64.StdEncoding.DecodeString(base64EncodedData); err != nil {
		return &econet.DataFrame{}, err
	}

	var data = econet.DataFrame{}
	data.DstStn = decodedFrame[0]
	data.DstNet = decodedFrame[1]
	data.SrcStn = decodedFrame[2]
	data.SrcNet = decodedFrame[3]

	// Add any data (Scouts don't have data)
	if len(decodedFrame) > 4 {
		data.Data = decodedFrame[4:]
	}
	return &data, nil
}
