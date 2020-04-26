package main

import (
	"flag"
	"fmt"
	"github.com/nikitawootten/cmsc483-project/common"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"net/url"
	"time"
	"math/rand"
	"strconv"
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

func makeKnownToParent(address, parentAddress string) {
	addressUrl, err := url.Parse(address)
	if err != nil {
		log.Printf("Could not parse address: %s, no connection to parent\n", err.Error())
	}

	conn, err := websocket.Dial(parentAddress, "", "http://localhost")
	if err != nil {
		log.Printf("Could not connect to websocket: %s, no connection to parent\n", err.Error())
	}
	defer conn.Close()

	req := common.NewClientReq{Address: addressUrl}
	err = websocket.JSON.Send(conn, req)
	if err != nil {
		log.Printf("Could not send initial request: %s, killing connection\n", err.Error())
	}

	for {
		time.Sleep(time.Minute * 1)
		// todo send some metrics
	}
}

func main() {
	var parentAddress = flag.String("parentAddress", "0.0.0.0:8080", "parent to connect to (empty for no parent)")
	var address = flag.String("address", ":8081", "address of self")
	flag.Parse()

	if *parentAddress != "" {
		go makeKnownToParent("http://0.0.0.0:8081", "ws://"+*parentAddress+"/client")
	}

	http.HandleFunc("/hello", helloWorld)

	err := http.ListenAndServe(*address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
