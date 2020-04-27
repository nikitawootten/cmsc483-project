package common

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
)

// see https://stackoverflow.com/questions/28322997/how-to-get-a-list-of-values-into-a-flag-in-golang
type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func ParseFlagsClient() (NewClientReq, []string, string, error) {
	var parentLBs = arrayFlags{}
	flag.Var(&parentLBs, "parentLB", "Parent load balancers to attempt to connect to, protocol and endpoint will be added automatically (ex. 0.0.0.0:8080)")
	var port = flag.Int("port", 8080, "Port of machine")

	// Needed to build NewClientReq
	var weight = flag.Int("weight", 1, "Weighted-round-robin weight definition to parent LB(s)")
	var callbackAddress = flag.String("callbackAddress", "", "URL that LB(s) will use to connect to self (Required)")

	flag.Parse()

	//log.Print("ParsedArgs: parentLBs=", parentLBs, " selfWebAddress=", *selfWebAddress, " callbackAddress=", *callbackAddress)

	req := NewClientReq{}

	selfWebAddress := fmt.Sprintf(":%d", *port)

	if len(parentLBs) > 0 {
		if *callbackAddress == "" {
			// attempt to read environment variable
			*callbackAddress = os.Getenv("IP_ADDR")
			if *callbackAddress == "" {
				return NewClientReq{}, nil, "", errors.New("callback address omitted")
			}
		}
		*callbackAddress = "http://" + *callbackAddress + selfWebAddress
		callbackAddressUrl, err := url.Parse(*callbackAddress)
		if err != nil {
			return NewClientReq{}, nil, "", err
		}

		req = NewClientReq{
			Address: callbackAddressUrl,
			Weight:  *weight,
		}
	}

	return req, parentLBs, selfWebAddress, nil
}

func ParseFlagsLB() (NewClientReq, []string, string, string, error) {
	var parentLBs = arrayFlags{}
	flag.Var(&parentLBs, "parentLB", "Parent load balancers to attempt to connect to, protocol and endpoint will be added automatically (ex. 0.0.0.0:8080)")
	var selfWebAddress = flag.String("webServerAddress", "", "Address used by web server config (specify 0.0.0.0:{some_port} to make globally accessible) (Required)")

	// Needed to build NewClientReq
	var weight = flag.Int("weight", 1, "Weighted-round-robin weight definition to parent LB(s)")
	var callbackAddress = flag.String("callbackAddress", "", "URL that LB(s) will use to connect to self (Required)")

	var algorithm = flag.String("algorithm", "round-robin", "Algorithm to use for scheduling")

	flag.Parse()

	if *selfWebAddress == "" {
		log.Println("Must specify webServerAddress")
		flag.PrintDefaults()
		os.Exit(1)
	}

	//log.Print("ParsedArgs: parentLBs=", parentLBs, " selfWebAddress=", *selfWebAddress, " callbackAddress=", *callbackAddress)

	req := NewClientReq{}

	if len(parentLBs) > 0 {
		callbackAddressUrl, err := url.Parse(*callbackAddress)
		if err != nil {
			return NewClientReq{}, nil, "", "", err
		}

		req = NewClientReq{
			Address: callbackAddressUrl,
			Weight:  *weight,
		}
	}

	return req, parentLBs, *selfWebAddress, *algorithm, nil
}
