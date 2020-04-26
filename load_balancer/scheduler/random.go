package scheduler

import (
	"math/rand"
	"net/http"
)

type RandomScheduler struct {
	clients []*Client
}

func NewRandomScheduler() *RandomScheduler {
	return &RandomScheduler{}
}

func (r *RandomScheduler) NewClient(client *Client) error {
	r.clients = append(r.clients, client)
	return nil
}

func (r *RandomScheduler) GetNext(_ *http.Request) (*Client, error) {
	if len(r.clients) == 0 {
		return nil, ErrNoClients
	}

	return r.clients[rand.Intn(len(r.clients))], nil
}
