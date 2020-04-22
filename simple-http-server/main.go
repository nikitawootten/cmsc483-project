package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nikitawootten/cmsc483-project/common"
	"log"
	"net/http"
	"net/url"
)

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "hi there!")
	if err != nil {
		log.Fatal(err)
	}
}

func makeKnownToParent(address, parentAddress string) error {
	addressUrl, err := url.Parse(address)
	if err != nil {
		return err
	}
	req := common.NewClientReq{Address: addressUrl}
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return err
	}
	_, err = http.Post(parentAddress, "application/json", bytes.NewBuffer(reqJSON))

	return err
}

func main() {
	var parentAddress = flag.String("parentAddress", "http://0.0.0.0:8080", "parent to connect to (empty for no parent)")
	var address = flag.String("address", ":8081", "address of self")
	flag.Parse()

	if *parentAddress != "" {
		err := makeKnownToParent("http://0.0.0.0:8081/client", *parentAddress)
		if err != nil {
			log.Fatal(err)
		}
	}

	http.HandleFunc("/hello", helloWorld)

	err := http.ListenAndServe(*address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
