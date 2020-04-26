package main

import (
	"fmt"
	"github.com/nikitawootten/cmsc483-project/common"
	"log"
	"net/http"
)

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	log.Println("New request!")
	_, err := fmt.Fprintf(w, "hi there!")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	req, lbs, address, err := common.ParseFlagsClient()
	if err != nil {
		log.Fatal("Failed to parse args:", err)
	}
	common.ConnectToParentLBs(req, lbs)

	http.HandleFunc("/hello", helloWorld)

	log.Println("Mapped routes, listening on ", address)

	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
