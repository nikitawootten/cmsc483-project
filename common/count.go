package common

import (
	"net/http"
	"sync/atomic"
)

type ConnectionCounter struct {
	count *uint32
}

func NewConnectionCounterFromHeartbeat(heartbeat *ClientHeartbeat) *ConnectionCounter {
	return &ConnectionCounter{count: &heartbeat.Connections}
}

// WrapHttp wraps a http.Handler counting the number of connections active at any time
func (cc *ConnectionCounter) WrapHttp(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint32(cc.count, 1)
		// defer ensures that it runs even during a panic
		defer atomic.AddUint32(cc.count, ^uint32(0)) // decrement uint32 (see https://golang.org/pkg/sync/atomic/#AddUint64)

		handler(w, r)
	}
}

func (cc ConnectionCounter) Get() uint32 {
	return *cc.count
}
