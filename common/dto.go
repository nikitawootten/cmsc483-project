package common

import (
	"net/url"
	"sync/atomic"
)

type NewClientReq struct {
	Address *url.URL `json:"address"`
	Weight  int      `json:"weight,omitempty"` // only used for weighted round robin
}

type ClientHeartbeat struct {
	Connections uint32 `json:"connections"`
}

// todo update to hold heartbeat metrics
func (ch *ClientHeartbeat) Update(newHeartbeat *ClientHeartbeat) {
	atomic.StoreUint32(&ch.Connections, newHeartbeat.Connections)
}
