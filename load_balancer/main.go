package main

import (
	"flag"
	"github.com/nikitawootten/cmsc483-project/load_balancer/scheduler"
	"github.com/nikitawootten/cmsc483-project/load_balancer/service"
	"log"
	"net/http"
)

func main() {
	// parse flags
	var parentAddress = flag.String("parentAddress", "", "parent to connect to (empty for no parent)")
	var algorithm = flag.String("algorithm", scheduler.RoundRobin, "which scheduling algorithm to use")
	var address = flag.String("address", ":8080", "address of self")
	flag.Parse()

	alg, err := scheduler.GetSchedulerByName(*algorithm)
	if err != nil {
		log.Fatal(err)
	}
	lb := service.NewLoadBalancer(*parentAddress, alg)

	// the parent communication system (register a client, get list of active clients)
	http.Handle("/client", lb.BuildClientHandlerFunc())

	// the load balancer itself
	http.HandleFunc("/", lb.BuildNewConnectionFunc())

	log.Println("Mapped routes, listening on ", *address)

	err = http.ListenAndServe(*address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
