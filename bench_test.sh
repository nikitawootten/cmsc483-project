#!/bin/bash
# replace 0.0.0.0 with IP address server is on

#$1 = address, $2 = output file prefix, $3 = image file name

ADDR1=$1
ADDR1+="/hello_world"
PREFIX=$2
echo "Stress tests for Load-balancing"
echo "Test hello_world endpoint"
ab -g ${PREFIX}-hello1.tsv -n 10000 -c 200 "$ADDR1"
sleep 10
ab -g ${PREFIX}-hello2.tsv -n 50000 -c 200 "$ADDR1" 
sleep 10
ab -g ${PREFIX}-hello3.tsv -n 100000 -c 200 "$ADDR1"
sleep 10
wrk -t 4 -c101 -d30s "$ADDR1" >> ${PREFIX}-hello1.txt
sleep 10
wrk -t 4 -c202 -d30s "$ADDR1" >> ${PREFIX}-hello2.txt
sleep 10
wrk -t 4 -c303 -d30s "$ADDR1" >> ${PREFIX}-hello3.txt
sleep 10


ADDR2=$1
ADDR2+="/fib"
echo "Test fibonacci endpoint" 
ab -g ${PREFIX}-fib1.tsv -n 10000 -c 200 "$ADDR2"  
sleep 10
ab -g ${PREFIX}-fib2.tsv -n 50000 -c 200 "$ADDR2" 
sleep 10
ab -g ${PREFIX}-fib3.tsv -n 100000 -c 200 "$ADDR2" 
sleep 10
wrk -t 4 -c101 -d30s "$ADDR2" >> ${PREFIX}-fib1.txt
sleep 10
wrk -t 4 -c202 -d30s "$ADDR2" >> ${PREFIX}-fib2.txt
sleep 10
wrk -t 4 -c303 -d30s "$ADDR2" >> ${PREFIX}-fib3.txt
sleep 10


ADDR3=$1
ADDR3+="/resize"
echo "Test resize endpoint" 

ab -g ${PREFIX}-resize1.tsv -n 1000 -c 100 -p $3 -T image/jpeg "$ADDR3" 
sleep 10
ab -g ${PREFIX}-resize2.tsv -n 5000 -c 100 -p $3 -T image/jpeg "$ADDR3" 
sleep 10
ab -g ${PREFIX}-resize3.tsv -n 10000 -c 100 -p $3 -T image/jpeg "$ADDR3" 
sleep 10
wrk -s postimg.lua -t 4 -c 10 -d 10s "$ADDR3" >> ${PREFIX}-resize1.txt
sleep 10
wrk -s postimg.lua -t 4 -c 20 -d 10s "$ADDR3" >> ${PREFIX}-resize2.txt
sleep 10
wrk -s postimg.lua -t 4 -c 30 -d 10s "$ADDR3" >> ${PREFIX}-resize3.txt



