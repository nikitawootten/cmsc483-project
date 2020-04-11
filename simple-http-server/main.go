package main

import (
	"fmt"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/log"
	"net/http"
)

func helloWorld(w http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprintf(w, "hi there!")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/hello_world", helloWorld)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
