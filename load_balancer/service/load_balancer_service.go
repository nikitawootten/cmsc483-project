package service

import (
	"github.com/nikitawootten/cmsc483-project/common"
	"github.com/nikitawootten/cmsc483-project/load_balancer/scheduler"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"time"
)

type LoadBalancer struct {
	parentAddr string
	scheduler  scheduler.IScheduler
}

func NewLoadBalancer(algorithm scheduler.IScheduler) LoadBalancer {
	return LoadBalancer{
		scheduler: algorithm,
	}
}

func (lb *LoadBalancer) BuildClientHandlerFunc() websocket.Handler {
	return func(ws *websocket.Conn) {
		log.Printf("Client connected from %s, receiving clientReq data...\n", ws.RemoteAddr())

		var clientReq common.NewClientReq
		err := websocket.JSON.Receive(ws, &clientReq)
		if err != nil {
			log.Printf("Error: Malformed clientReq handshake: %s, killing connection\n", err.Error())
			return
		}

		//log.Print("addr: ", ws.Config().Origin)

		client := scheduler.NewClient(clientReq)
		err = lb.scheduler.NewClient(&client)
		if err != nil {
			log.Printf("Error: Could not add clientReq: %s, killing connection\n", err.Error())
		}

		defer func() {
			client.Active = false
			log.Printf("Disconnecting from client at address: %s", clientReq.Address)
			err = ws.Close()
			if err != nil {
				log.Printf("Error: Could not close connection: %s, conneciton killed\n", err.Error())
			}
		}()

		errCount := 0
		for {
			var heartbeat common.ClientHeartbeat
			err := websocket.JSON.Receive(ws, &heartbeat)
			if err != nil {
				log.Printf("Warning: Malformed heartbeat: %s, continuing\n", err.Error())
				errCount++
				if errCount > 3 {
					return // kill connection
				}
				time.Sleep(time.Second)
				continue
			}
			errCount = 0 // reset error count

			client.Heartbeat.Update(&heartbeat)
		}
	}
}

func (lb *LoadBalancer) BuildNewConnectionFunc() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("New Connection!")
		client, err := lb.scheduler.GetNext(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			client.Proxy.ServeHTTP(w, r)
		}
	}
}
