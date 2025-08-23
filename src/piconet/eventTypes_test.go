package piconet

import (
	"encoding/base64"
	"reflect"
	"testing"

	"github.com/johnnewcombe/econet-simple-server/src/econet"
)

// helper to build and encode frames
func b64ScoutFrame(s econet.ScoutFrame) string {
	raw := append([]byte{
		s.DstStn,
		s.DstNet,
		s.SrcStn,
		s.SrcNet,
		s.ControlByte,
		s.Port,
	}, s.Data...)
	return base64.StdEncoding.EncodeToString(raw)
}

func b64DataFrame(d econet.DataFrame) string {
	raw := append([]byte{
		d.DstStn,
		d.DstNet,
		d.SrcStn,
		d.SrcNet,
		//d.ReplyPort,
		//d.FunctionCode,
	}, d.Data...)
	return base64.StdEncoding.EncodeToString(raw)
}

func TestNewRxTransmit_Success(t *testing.T) {

	scout := econet.ScoutFrame{
		NetHeader:   econet.NetHeader{DstStn: 0x11, DstNet: 0x22, SrcStn: 0x33, SrcNet: 0x44},
		ControlByte: 0x55,
		Port:        0x66,
		Data:        []byte{},
	}

	data := econet.DataFrame{
		NetHeader: econet.NetHeader{DstStn: 0x10, DstNet: 0x20, SrcStn: 0x30, SrcNet: 0x40},
		//ReplyPort:    0x50,
		//FunctionCode: 0x60,
		Data: []byte("HELLO"),
	}

	args := []string{b64ScoutFrame(scout), b64DataFrame(data)}

	rxt, err := NewRxTransmit(args)
	if err != nil {
		t.Fatalf("NewRxTransmit returned error: %v", err)
	}
	if rxt == nil {
		t.Fatalf("NewRxTransmit returned nil RxTransmit")
	}

	// Verify scout fields
	if rxt.ScoutFrame == nil {
		t.Fatalf("ScoutFrame is nil")
	}
	if rxt.ScoutFrame.DstStn != scout.DstStn || rxt.ScoutFrame.DstNet != scout.DstNet ||
		rxt.ScoutFrame.SrcStn != scout.SrcStn || rxt.ScoutFrame.SrcNet != scout.SrcNet ||
		rxt.ScoutFrame.ControlByte != scout.ControlByte || rxt.ScoutFrame.Port != scout.Port {
		t.Errorf("ScoutFrame fields mismatch: got %+v want %+v", *rxt.ScoutFrame, scout)
	}
	if !reflect.DeepEqual(rxt.ScoutFrame.Data, scout.Data) {
		t.Errorf("ScoutFrame.Data mismatch: got % X want % X", rxt.ScoutFrame.Data, scout.Data)
	}

	// Verify data fields
	if rxt.DataFrame == nil {
		t.Fatalf("DataFrame is nil")
	}
	if rxt.DataFrame.DstStn != data.DstStn || rxt.DataFrame.DstNet != data.DstNet ||
		rxt.DataFrame.SrcStn != data.SrcStn || rxt.DataFrame.SrcNet != data.SrcNet { //} ||
		//		rxt.DataFrame.ReplyPort != data.ReplyPort || rxt.DataFrame.FunctionCode != data.FunctionCode {
		t.Errorf("DataFrame fields mismatch: got %+v want %+v", *rxt.DataFrame, data)
	}
	if !reflect.DeepEqual(rxt.DataFrame.Data, data.Data) {
		t.Errorf("DataFrame.Data mismatch: got % X want % X", rxt.DataFrame.Data, data.Data)
	}
}

func TestNewRxTransmit_Errors(t *testing.T) {
	// Valid frames for one side
	validScout := b64ScoutFrame(econet.ScoutFrame{
		NetHeader:   econet.NetHeader{DstStn: 1, DstNet: 2, SrcStn: 3, SrcNet: 4},
		ControlByte: 5,
		Port:        6,
	})
	validData := b64DataFrame(econet.DataFrame{
		NetHeader: econet.NetHeader{DstStn: 7, DstNet: 8, SrcStn: 9, SrcNet: 10},
		//ReplyPort:    11,
		//FunctionCode: 12,
	})

	// Case 1: invalid base64 for scout
	if _, err := NewRxTransmit([]string{"!!!not-base64!!!", validData}); err == nil {
		t.Errorf("expected error for invalid scout base64, got nil")
	}

	// Case 2: invalid base64 for data
	if _, err := NewRxTransmit([]string{validScout, "***bad***"}); err == nil {
		t.Errorf("expected error for invalid data base64, got nil")
	}
}
