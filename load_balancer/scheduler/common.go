package scheduler

import (
	"errors"
	"github.com/nikitawootten/cmsc483-project/common"
	"net/http"
	"net/http/httputil"
)

type IScheduler interface {
	NewClient(client common.NewClientReq) error
	GetNext(r *http.Request) (*httputil.ReverseProxy, error)
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
