package scheduler

import (
	"net/http"
	"sync/atomic"
)

type RoundRobinScheduler struct {
	clients []*Client
	count   uint32 // atomic connection count
}

func NewRoundRobinScheduler() *RoundRobinScheduler {
	return &RoundRobinScheduler{}
}

func (rr *RoundRobinScheduler) NewClient(client *Client) error {
	for i, currClient := range rr.clients {
		if currClient.Init.Address == client.Init.Address {
			// existing client found, replace with new client
			rr.clients[i] = client
			return nil
		}
	}
	rr.clients = append(rr.clients, client)
	return nil
}

// GetNext does not use any information about the request, so the input is ignored
func (rr *RoundRobinScheduler) GetNext(_ *http.Request) (*Client, error) {
	if len(rr.clients) == 0 {
		return nil, ErrNoClients
	}

	count := atomic.AddUint32(&rr.count, 1)
	index := int(count) % len(rr.clients)
	ret := rr.clients[index]
	//log.Println(index, count)
	return ret, nil
}
