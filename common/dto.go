package common

import "net/url"

type NewClientReq struct {
	Address *url.URL `json:"address"`
	Weight  int      `json:"weight,omitempty"` // only used for weighted round robin
}

type ClientHeartbeat struct {
}
