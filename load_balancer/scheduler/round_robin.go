package scheduler

import (
	"errors"
	"github.com/nikitawootten/cmsc483-project/common"
	"net/http"
	"net/http/httputil"
	"sync"
)

type RoundRobinScheduler struct {
	lock    sync.RWMutex // todo implement locking
	proxies []*httputil.ReverseProxy
	curr    int
}

func NewRoundRobinScheduler() *RoundRobinScheduler {
	return &RoundRobinScheduler{}
}

func (rr *RoundRobinScheduler) NewClient(client common.NewClientReq) error {
	rp := httputil.NewSingleHostReverseProxy(client.Address)
	rr.proxies = append(rr.proxies, rp)
	return nil
}

// GetNext does not use any information about the request, so the input is ignored
func (rr *RoundRobinScheduler) GetNext(_ *http.Request) (*httputil.ReverseProxy, error) {
	ret := rr.proxies[rr.curr]
	rr.curr += 1
	return ret, errors.New("not implemented")
}
