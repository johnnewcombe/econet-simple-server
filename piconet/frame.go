package piconet

import (
	"encoding/base64"
	"fmt"
	"strings"
)

/*
+------+------+-----+-----+---------+------+-----------------------------+
| Dest | Dest | Src | Src | Control | Port |         Data                |
| Stn  | Net  | Stn | Net |  Byte   |      |                             |
+------+------+-----+-----+---------+------+-----------------------------+

	<-------- - - Packet Header - - ---------> <--- - - Packet Data - - --->

Consists of 6 bytes+data, A Scout is the same but without the data section.
*/
type Header struct {
	DstStn      byte
	DstNet      byte
	SrcStn      byte
	SrcNet      byte
	ControlByte byte
	Port        EconetPort
}

type Frame struct {
	ControlByte byte
	Port        EconetPort
	Data        []byte
	Header
}

// String returns a structured representation of the frame.
func (f *Frame) String() string {
	var sb = strings.Builder{}

	sb.WriteString(fmt.Sprintf("dst-stn=%02X, dst-net=%02X, src-stn=%02X, scr-net=%02X, ctrl-byte=%02X, port=%02X, port-description=%s",
		f.DstStn, f.DstNet, f.SrcStn, f.SrcNet, f.ControlByte, f.Port.Value, f.Port.Description))
	if len(f.Data) > 0 {
		sb.WriteString(fmt.Sprintf(", data=%02X", f.Data))
	}

	return sb.String()
}

// NewData Frame creates a nee instance of a Piconet data frame
func NewDataFrame(base64EncodedData string) (Frame, error) {
	var (
		port         EconetPort
		decodedFrame []byte
		err          error
	)

	if decodedFrame, err = base64.StdEncoding.DecodeString(base64EncodedData); err != nil {
		return Frame{}, err
	}

	var f = Frame{}
	f.DstStn = decodedFrame[0]
	f.DstNet = decodedFrame[1]
	f.SrcStn = decodedFrame[2]
	f.SrcNet = decodedFrame[3]
	f.ControlByte = decodedFrame[4]

	// Add the port
	if port, err = NewPort(decodedFrame[5]); err != nil {
		return Frame{}, err
	}
	f.Port = port

	// Add any data (Scouts don't have data)
	if len(decodedFrame) > 6 {
		f.Data = decodedFrame[6:]
	}

	return f, nil
}

/*
type ScoutFrame struct {
	Data []byte
	Header
}

func NewScoutFrame(base64EncodedData string) (Frame, error) {
	var (
		decodedFrame []byte
		err          error
	)

	if decodedFrame, err = base64.StdEncoding.DecodeString(base64EncodedData); err != nil {
		return Frame{}, err
	}

	var f = Frame{}
	f.DstStn = decodedFrame[0]
	f.DstNet = decodedFrame[1]
	f.SrcStn = decodedFrame[2]
	f.SrcNet = decodedFrame[3]
	f.ControlByte = decodedFrame[4]

	// Add the port
	if port, err = NewPort(decodedFrame[5]); err != nil {
		return Frame{}, err
	}
	f.Port = port

	// Notify has 'Scout Extra Data'

	if len(decodedFrame) > 4 {
		f.Data = decodedFrame[4:]
	}

	return f, nil
}

// String returns a structured representation of the frame.
func (f *ScoutFrame) String() string {
	var sb = strings.Builder{}
	sb.WriteString(fmt.Sprintf("dst-stn=%02X, dst-net=%02X, src-stn=%02X, scr-net=%02X, ctrl-byte=%02X, port=%02X, port-description=%s",
		f.DstStn, f.DstNet, f.SrcStn, f.SrcNet, f.ControlByte, f.Port.Value, f.Port.Description))
	return sb.String()
}
*/
