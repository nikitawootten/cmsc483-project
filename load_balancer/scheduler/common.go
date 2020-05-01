package scheduler

import (
	"errors"
	"github.com/nikitawootten/cmsc483-project/common"
	"net/http"
	"net/http/httputil"
)

type Client struct {
	Active    bool
	Proxy     httputil.ReverseProxy
	Init      common.NewClientReq
	Heartbeat common.ClientHeartbeat
}

func NewClient(init common.NewClientReq) Client {
	rp := httputil.NewSingleHostReverseProxy(init.Address)
	return Client{
		Active: true,
		Proxy:  *rp,
		Init:   init,
	}
}

type IScheduler interface {
	NewClient(client *Client) error
	GetNext(r *http.Request) (*Client, error)
}

const (
	RoundRobin       = "round-robin"
	Random           = "random"
	LeastConnections = "least-connections"
)

var ErrNoClients = errors.New("no clients to load balance")
var ErrUnkSchedAlg = errors.New("unknown scheduling algorithm")

func GetSchedulerByName(algorithm string) (IScheduler, error) {
	switch algorithm {
	case RoundRobin:
		return NewRoundRobinScheduler(), nil
	case Random:
		return NewRandomScheduler(), nil
	case LeastConnections:
		return NewLeastConnectionsScheduler(), nil
	default:
		return nil, ErrUnkSchedAlg
	}
}
