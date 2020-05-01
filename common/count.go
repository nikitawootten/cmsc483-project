package common

import (
	"net/http"
	"sync/atomic"
)

type ConnectionCounter struct {
	count uint32
}

func NewConnectionCounter() *ConnectionCounter {
	return &ConnectionCounter{count: 0}
}

// WrapHttp wraps a http.Handler counting the number of connections active at any time
func (cc *ConnectionCounter) WrapHttp(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint32(&cc.count, 1)
		// defer ensures that it runs even during a panic
		defer atomic.AddUint32(&cc.count, -1)

		handler(w, r)
	}
}

func (cc ConnectionCounter) Get() uint32 {
	return cc.count
}
