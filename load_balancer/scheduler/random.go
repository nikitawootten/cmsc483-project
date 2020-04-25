package scheduler

import (
	"github.com/nikitawootten/cmsc483-project/common"
	"math/rand"
	"net/http"
	"net/http/httputil"
)

type RandomScheduler struct {
	proxies []*httputil.ReverseProxy
}

func NewRandomScheduler() *RandomScheduler {
	return &RandomScheduler{}
}

func (r *RandomScheduler) NewClient(client common.NewClientReq) error {
	rp := httputil.NewSingleHostReverseProxy(client.Address)
	r.proxies = append(r.proxies, rp)
	return nil
}

func (r *RandomScheduler) GetNext(_ *http.Request) (*httputil.ReverseProxy, error) {
	if len(r.proxies) == 0 {
		return nil, ErrNoClients
	}

	return r.proxies[rand.Intn(len(r.proxies))], nil
}
