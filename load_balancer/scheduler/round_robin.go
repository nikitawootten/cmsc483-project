package scheduler

import (
	"github.com/nikitawootten/cmsc483-project/common"
	"log"
	"net/http"
	"net/http/httputil"
	"sync/atomic"
)

type RoundRobinScheduler struct {
	proxies []*httputil.ReverseProxy
	count   uint32 // atomic connection count
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
	if len(rr.proxies) == 0 {
		return nil, ErrNoClients
	}

	count := atomic.AddUint32(&rr.count, 1)
	index := int(count) % len(rr.proxies)
	ret := rr.proxies[index]
	log.Println(index, count)
	return ret, nil
}
