package piconet

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/johnnewcombe/econet-simple-server/src/econet"
)

type Event struct {
	Cmd     string
	Args    []string
	CmdText string
}

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

type StatusResponse struct {
	MajorVersion int
	MinorVersion int
	Patch        int
	Station      int
	StatusReg    int
	Mode         int
}

func NewStatusResponse(statusText []string) (*StatusResponse, error) {

	var (
		err     error
		sr      = StatusResponse{}
		version []string
	)

	if len(statusText) != 4 {
		return &StatusResponse{}, errors.New("invalid status response")
	}

	version = strings.Split(statusText[0], ".")

	if len(version) != 3 {
		return &StatusResponse{}, errors.New("invalid status response")
	}

	if sr.MajorVersion, err = strconv.Atoi(version[0]); err != nil {
		return &StatusResponse{}, errors.New("invalid staus response (major version)")
	}
	if sr.MinorVersion, err = strconv.Atoi(version[1]); err != nil {
		return &StatusResponse{}, errors.New("invalid staus response (minor version)")
	}
	if sr.Patch, err = strconv.Atoi(version[2]); err != nil {
		return &StatusResponse{}, errors.New("invalid staus response (patch version)")
	}
	if sr.Station, err = strconv.Atoi(statusText[1]); err != nil || sr.Station > 254 || sr.Station < 1 {
		return &StatusResponse{}, errors.New("invalid staus response (station number)")
	}
	if sr.StatusReg, err = strconv.Atoi(statusText[2]); err != nil || sr.StatusReg > 255 || sr.StatusReg < 0 {
		return &StatusResponse{}, errors.New("invalid staus response (status register)")
	}
	if sr.Mode, err = strconv.Atoi(statusText[3]); err != nil || sr.Mode < 0 || sr.Mode > 3 {
		return &StatusResponse{}, errors.New("invalid staus response (mode)")
	}

	return &sr, nil
}

func (s *StatusResponse) String() string {

	var sb = strings.Builder{}
	sb.WriteString(fmt.Sprintf("maj-ver=%d, min-ver=%d, patch=%d, stn=%d, status-reg=%02Xh, mode=%d, mode-desc=%s",
		s.MajorVersion, s.MinorVersion, s.Patch, s.Station, s.StatusReg, s.Mode, ModeMap[s.Mode]))

	return sb.String()
}

var ModeMap = map[int]string{
	0: "STOP",
	1: "LISTEN",
	2: "MONITOR",
}

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
		rxt.ScoutFrame.DstStn, rxt.ScoutFrame.DstNet, rxt.ScoutFrame.SrcStn, rxt.ScoutFrame.SrcNet,
		rxt.ScoutFrame.ControlByte, rxt.ScoutFrame.Port, econet.PortMap[rxt.ScoutFrame.Port]))
	sb.WriteString(", ")
	sb.WriteString(fmt.Sprintf("data-dst-stn=%02X, data-dst-net=%02X, data-src-stn=%02X, data-scr-net=%02X, reply-port=%02X, function-code=%02x",
		rxt.DataFrame.DstStn, rxt.DataFrame.DstNet, rxt.DataFrame.SrcStn, rxt.DataFrame.SrcNet,
		rxt.DataFrame.ReplyPort, rxt.DataFrame.FunctionCode))

	if len(rxt.DataFrame.Data) > 0 {
		sb.WriteString(fmt.Sprintf(", data-bytes=[% 02X]", rxt.DataFrame.Data))
	}

	return sb.String()
}

//func (rxt *RxTransmit) Command() string {

//	if len(rxt.DataFrame.Data) > 0 {
//		return strings.TrimRight(string(rxt.DataFrame.Data), "\r")
//	}
//	return ""
//}

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
	if len(decodedFrame) > 6 {
		scout.Data = decodedFrame[6:]
	}

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
	data.ReplyPort = decodedFrame[4]
	data.FunctionCode = decodedFrame[5]

	// Add any data (Scouts don't have data)
	if len(decodedFrame) > 6 {
		data.Data = decodedFrame[6:]
	}
	return &data, nil
}

/*
type TxResult struct {
	Result string
	Ok     bool
}

func NewTxResult(result string) TxResult {

	txResult := TxResult{
		Result: result,
	}
	if result == "OK" {
		txResult.Ok = true
	}
	return txResult
}
*/
