#!/usr/bin/env bash
# Starts a sample configuration with 3 web servers and 2 lbs

go run ./load_balancer/main.go -port 8080 &
go run ./load_balancer/main.go -port 8081 -callbackAddress 0.0.0.0:8081 -parentLB 0.0.0.0:8080 &
# http server connected to both
go run ./simple-http-server/main.go -port 8082 -callbackAddress 0.0.0.0:8082 -parentLB 0.0.0.0:8080 -parentLB 0.0.0.0:8081 &
# http servers connected to one or the other
go run ./simple-http-server/main.go -port 8083 -callbackAddress 0.0.0.0:8083 -parentLB 0.0.0.0:8080 &
go run ./simple-http-server/main.go -port 8084 -callbackAddress 0.0.0.0:8084 -parentLB 0.0.0.0:8081
wait