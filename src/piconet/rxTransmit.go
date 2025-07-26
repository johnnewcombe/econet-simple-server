package piconet

import (
	"encoding/base64"
	"fmt"
	"strings"
)

type DataFrame struct {
	Header
	Data []byte
}

type RxTransmit struct {
	ScoutFrame ScoutFrame
	DataFrame  DataFrame
}

func NewRxTransmit(eventArgs []string) (RxTransmit, error) {

	var (
		err   error
		rxt   RxTransmit
		scout ScoutFrame
		data  DataFrame
	)

	if scout, err = newScoutFrame(eventArgs[0]); err != nil {
		return RxTransmit{}, err
	}

	if data, err = newDataFrame(eventArgs[1]); err != nil {
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
	sb.WriteString(fmt.Sprintf("data-dst-stn=%02X, data-dst-net=%02X, data-src-stn=%02X, data-scr-net=%02X, data-ctrl-byte=%02X",
		rxt.DataFrame.DstStn, rxt.DataFrame.DstNet, rxt.DataFrame.SrcStn, rxt.DataFrame.SrcNet))

	if len(rxt.DataFrame.Data) > 0 {
		sb.WriteString(fmt.Sprintf(", data-bytes=[% 02X]", rxt.DataFrame.Data))
	}

	return sb.String()
}

func (rxt *RxTransmit) Command() string {

	if len(rxt.DataFrame.Data) > 4 {
		return string(rxt.DataFrame.Data[4:])
	}
	return ""
}

func newDataFrame(base64EncodedData string) (DataFrame, error) {

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

	// Add any data (Scouts don't have data)
	if len(decodedFrame) > 4 {
		data.Data = decodedFrame[4:]
	}
	return data, nil
}
