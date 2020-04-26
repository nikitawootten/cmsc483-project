package scheduler

import (
	"errors"
	"github.com/nikitawootten/cmsc483-project/common"
	"net/http"
	"net/http/httputil"
)

type Client struct {
	Proxy     httputil.ReverseProxy
	Init      common.NewClientReq
	Heartbeat common.ClientHeartbeat
}

func NewClient(init common.NewClientReq) Client {
	rp := httputil.NewSingleHostReverseProxy(init.Address)
	return Client{
		Proxy: *rp,
		Init:  init,
	}
}

type IScheduler interface {
	NewClient(client *Client) error
	GetNext(r *http.Request) (*Client, error)
}

const (
	RoundRobin = "round-robin"
	Random     = "random"
)

var ErrNoClients = errors.New("no clients to load balance")

func GetSchedulerByName(algorithm string) (IScheduler, error) {
	switch algorithm {
	case RoundRobin:
		return NewRoundRobinScheduler(), nil
	case Random:
		return NewRandomScheduler(), nil
	default:
		return nil, errors.New("unknown scheduling algorithm")
	}
}
