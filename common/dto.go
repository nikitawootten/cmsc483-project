package common

import "net/url"

type NewClientReq struct {
	Address *url.URL `json:"address"`
}
