package piconet

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type StatusResponse struct {
	MajorVersion int
	MinorVersion int
	Patch        int
	Station      int
	StatusReg    int
	Mode         int
}

func NewStatusResponse(statusText []string) (StatusResponse, error) {

	var err error

	if len(statusText) != 4 {
		return StatusResponse{}, errors.New("invalid status response")
	}

	version := strings.Split(statusText[0], ".")

	if len(version) != 3 {
		return StatusResponse{}, errors.New("invalid status response")
	}
	var sr = StatusResponse{}

	if sr.MajorVersion, err = strconv.Atoi(version[0]); err != nil {
		return StatusResponse{}, errors.New("invalid staus response (major version)")
	}
	if sr.MinorVersion, err = strconv.Atoi(version[1]); err != nil {
		return StatusResponse{}, errors.New("invalid staus response (minor version)")
	}
	if sr.Patch, err = strconv.Atoi(version[2]); err != nil {
		return StatusResponse{}, errors.New("invalid staus response (patch version)")
	}
	if sr.Station, err = strconv.Atoi(statusText[1]); err != nil || sr.Station > 254 || sr.Station < 1 {
		return StatusResponse{}, errors.New("invalid staus response (station number)")
	}
	if sr.StatusReg, err = strconv.Atoi(statusText[2]); err != nil || sr.StatusReg > 255 || sr.StatusReg < 0 {
		return StatusResponse{}, errors.New("invalid staus response (status register)")
	}
	if sr.Mode, err = strconv.Atoi(statusText[3]); err != nil || sr.Mode < 0 || sr.Mode > 3 {
		return StatusResponse{}, errors.New("invalid staus response (mode)")
	}

	return sr, nil
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
