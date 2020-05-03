#!/bin/bash
# replace 0.0.0.0 with IP address server is on

#$1 = address, $2 = output file prefix, $3 = image file name
PREFIX=$2
ADDR=$1

FIB="$ADDR/fib"
HELLO="$ADDR/hello_world"
RESIZE="$ADDR/resize"

echo "Stress tests for Load-balancing"
echo "Test hello_world endpoint"
ab -g ${PREFIX}-hello.tsv -n 100000 -c 100 "$HELLO"
sleep 5
echo "Test hello_world endpoint wrk"
wrk -t 4 -c100 -d30s "http://$HELLO" >> ${PREFIX}-hello.txt
sleep 5

echo "Test fibonacci endpoint" 
ab -g ${PREFIX}-fib.tsv -n 100000 -c 100 "$FIB"
sleep 5
echo "Test fibonacci endpoint wrk"
wrk -t 4 -c100 -d30s "http://$FIB" >> ${PREFIX}-fib.txt
sleep 5

echo "Test resize endpoint" 
ab -g ${PREFIX}-resize.tsv -n 5000 -c 100 -p $3 -T image/jpeg "$RESIZE"
sleep 5
echo "Test resize endpoint wrk"
wrk -s postimg.lua -t 4 -c 20 -d 30s "http://$RESIZE" >> ${PREFIX}-resize.txt
