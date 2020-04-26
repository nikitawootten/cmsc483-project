#!/usr/bin/env bash
# Starts a sample configuration with 3 web servers and 2 lbs

go run ./load_balancer/main.go -webServerAddress :8080 &
go run ./load_balancer/main.go -webServerAddress :8081 -callbackAddress http://0.0.0.0:8081 -parentLB 0.0.0.0:8080 &
# http server connected to both
go run ./simple-http-server/main.go -webServerAddress :8082 -callbackAddress http://0.0.0.0:8082 -parentLB 0.0.0.0:8080 -parentLB 0.0.0.0:8081 &
# http servers connected to one or the other
go run ./simple-http-server/main.go -webServerAddress :8083 -callbackAddress http://0.0.0.0:8083 -parentLB 0.0.0.0:8080 &
go run ./simple-http-server/main.go -webServerAddress :8084 -callbackAddress http://0.0.0.0:8084 -parentLB 0.0.0.0:8081
wait