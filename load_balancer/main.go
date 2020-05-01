package main

import (
	"github.com/nikitawootten/cmsc483-project/common"
	"github.com/nikitawootten/cmsc483-project/stats"
	"github.com/nikitawootten/cmsc483-project/load_balancer/scheduler"
	"github.com/nikitawootten/cmsc483-project/load_balancer/service"
	"log"
	"net/http"
)

func main() {
	req, lbs, address, algorithm, err := common.ParseFlags(true)
	if err != nil {
		log.Fatal("Failed to parse args:", err)
	}
	common.ConnectToParentLBs(req, lbs)

	alg, err := scheduler.GetSchedulerByName(algorithm)
	if err != nil {
		log.Fatal(err)
	}
	lb := service.NewLoadBalancer(alg)
	
	// the parent communication system (register a client, get list of active clients)
	http.Handle("/client", lb.BuildClientHandlerFunc())

	// the load balancer itself
	http.HandleFunc("/", lb.BuildNewConnectionFunc())

	log.Println("Mapped routes, listening on ", address)

	go stats.ExecuteCronJob()

	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
	}

}
