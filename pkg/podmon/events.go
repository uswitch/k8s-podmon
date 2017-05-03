package podmon

import (
	"context"
	"sync"
)

// Alerter defines the required methods to send alerts via podmon
type Alerter interface {
	EventLoop(context.Context, sync.WaitGroup, chan interface{})
	Send(interface{}) (interface{}, error)
}
