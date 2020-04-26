package main

import (
	"fmt"
	"github.com/nikitawootten/cmsc483-project/common"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func fib(n int) int {
	varOne := 0
	varTwo := 1
	for i := 0; i < n; i++ {
		temp := varOne
		varOne = varTwo
		varTwo = temp + varOne
	}
	return varOne
}

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	log.Println("New request!")
	_, err := fmt.Fprintf(w, "hi there!")

	fmt.Fprintf(w, "\n")
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < (rand.Intn(100 - 90) + 90); i++ {

		fmt.Fprintf(w, strconv.Itoa(fib(i)) + " ")
	}

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
