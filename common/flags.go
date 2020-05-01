package common

import (
	"flag"
	"fmt"
	"net/http"
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

func ParseFlags(isLB bool) (NewClientReq, []string, string, string, error) {
	var parentLBs = arrayFlags{}
	flag.Var(&parentLBs, "parentLB", "Parent load balancers to attempt to connect to, protocol and endpoint will be added automatically (ex. 0.0.0.0:8080)")
	var port = flag.Int("port", 8080, "Port of machine")

	// Needed to build NewClientReq
	var weight = flag.Int("weight", 1, "Weighted-round-robin weight definition to parent LB(s)")
	var callbackAddress = flag.String("callbackAddress", "", "URL that LB(s) will use to connect to self (Required)")

	var algorithm = ""
	if isLB {
		flag.StringVar(&algorithm, "algorithm", "round-robin", "Algorithm to use for scheduling")
	}

	var maxIdleCon = flag.Int("maxIdleCon", 100, "Maximum idle http connections per host")

	flag.Parse()

	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = *maxIdleCon

	req := NewClientReq{}

	selfWebAddress := fmt.Sprintf(":%d", *port)

	if len(parentLBs) > 0 {
		if *callbackAddress == "" {
			*callbackAddress = os.Getenv("IP_ADDR")
			*callbackAddress = fmt.Sprint(*callbackAddress, ":", *port)
		}

		callbackAddressUrl, err := url.Parse(fmt.Sprint("http://", *callbackAddress))
		if err != nil {
			return NewClientReq{}, nil, "", "", err
		}

		req = NewClientReq{
			Address: callbackAddressUrl,
			Weight:  *weight,
		}
	}

	return req, parentLBs, selfWebAddress, algorithm, nil
}
