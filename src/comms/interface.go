package comms

import (
	"context"
	"sync"
)

type funcDef func(bool, byte)

type CommunicationClient interface {
	Open(portName string) error
	Write(byt []byte) error
	Read(ctx context.Context, wg *sync.WaitGroup, c chan byte)
	Close() error
}
