package service

import (
	"encoding/json"
	"github.com/nikitawootten/cmsc483-project/common"
	"github.com/nikitawootten/cmsc483-project/load_balancer/scheduler"
	"log"
	"net/http"
)

type LoadBalancer struct {
	parentAddr string
	scheduler  scheduler.IScheduler
	clients    []string
}

func NewLoadBalancer(parentAddr string, algorithm scheduler.IScheduler) LoadBalancer {
	return LoadBalancer{
		parentAddr: parentAddr,
		scheduler:  algorithm,
		clients:    []string{},
	}
}

func (lb *LoadBalancer) BuildNewClientFunc() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("New Client!")
		var client common.NewClientReq

		err := json.NewDecoder(r.Body).Decode(&client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = lb.scheduler.NewClient(client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

func (lb *LoadBalancer) BuildNewConnectionFunc() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("New Connection!")
		rp, err := lb.scheduler.GetNext(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			rp.ServeHTTP(w, r)
		}
	}
}
