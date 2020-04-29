package common

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"time"
)

const maxConnectFails = 3
const reconnectTimeout = time.Second * 3

func ConnectToParentLBs(req NewClientReq, lbs []string) {
	for _, lb := range lbs {
		go MakeKnownToParent(req, lb)
	}
}

func MakeKnownToParent(req NewClientReq, parentAddress string) {
	origin := fmt.Sprint(req.Address.String())
	parentAddress = fmt.Sprint("ws://", parentAddress, "/client")

	connectFails := 0
	for {
		if connectFails > maxConnectFails {
			log.Printf("Maximum retries reached, cancelling connection attempts to parent")
			return
		}

		conn, err := websocket.Dial(parentAddress, "", origin)
		if err != nil {
			log.Printf("Could not connect to websocket: %s, no connection to parent\n", err.Error())
			connectFails += 1
			time.Sleep(reconnectTimeout)
			continue
		}


		log.Printf("Connected to parent lb %s, streaming metrics\n", parentAddress)

		err = websocket.JSON.Send(conn, req)
		if err != nil {
			log.Printf("Could not send initial request: %s, killing connection\n", err.Error())
			connectFails += 1
			time.Sleep(reconnectTimeout)
			continue
		}

		connectFails = 0

		for {
			time.Sleep(time.Minute * 1)
			// todo send some metrics
		}
		//log.Printf("Disconnected from lb %s, sleeping before attempting to reconnect!\n", parentAddress)
		//time.Sleep(reconnectTimeout)
	}
}
