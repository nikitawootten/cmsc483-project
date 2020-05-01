package scheduler

import (
	"math"
	"net/http"
	"sync/atomic"
)

type LeastConnectionsScheduler struct {
	clients []*Client
}

func NewLeastConnectionsScheduler() *LeastConnectionsScheduler {
	return &LeastConnectionsScheduler{}
}

func (lc *LeastConnectionsScheduler) NewClient(client *Client) error {
	for i, currClient := range lc.clients {
		if currClient.Init.Address == client.Init.Address {
			// existing client found, replace with new client
			lc.clients[i] = client
			return nil
		}
	}
	lc.clients = append(lc.clients, client)
	return nil
}

func (lc *LeastConnectionsScheduler) GetNext(_ *http.Request) (*Client, error) {
	var minCount uint32 = math.MaxUint32
	minIndex := -1
	for i, client := range lc.clients {
		count := client.Heartbeat.Connections
		if count <= minCount {
			minCount = count
			minIndex = i
		}
	}

	if minIndex == -1 {
		return nil, ErrNoClients
	}

	// update cached least connections value with new connection count
	atomic.AddUint32(&lc.clients[minIndex].Heartbeat.Connections, 1)

	return lc.clients[minIndex], nil
}
