#!/usr/bin/env bash
# This file performs load testing on a variable number of clients.
# Given a docker compose file, it will start all load balancers and progressively test the system as more CPUs are added
#
# Ensure that all clients have the form `client-{#}` and all load balancers contain the suffix `-lb` for this script to
# properly find all the relevant clients
# Also sure that the root load balancer is accessible from localhost and runs off port 8080

if [[ "$#" -ne 4 ]]; then
   echo "Usage: ./load-experiment.sh <compose-file> <min client num> <max client num> <output directory>"
   exit 2
fi

COMPOSE_FILE=$1 # The compose file to be tested
MIN_CLIENT_NUM=$2 # The point at which testing begins (used to prevent starving a load balancer)
MAX_CLIENT_NUM=$3 # Number of clients this experiment will test to
OUTPUT_ROOT=$4

LOAD_BALANCER_SUFFIX="lb"
SLEEP_INFRA_DURATION=5
SLEEP_TEST_DURATION=30 # amount of time to sleep between tests

# for some reason this is abstracted out?
IMG_RESIZE_PATH="test_img.jpeg"
ROOT_LB_ADDR="127.0.0.1:8080"

echo "========== Load Test =========="
echo "=== Starting load balancers ==="

docker-compose -f ${COMPOSE_FILE} config --services | grep ${LOAD_BALANCER_SUFFIX} | while read service ; do
    docker-compose -f ${COMPOSE_FILE} up -d --force-recreate ${service}
    echo "Started service ${service}..."
done

sleep ${SLEEP_INFRA_DURATION}
echo "=== Starting initial clients ==="

for (( i=1; i<=${MIN_CLIENT_NUM}; i++ ))
do
    docker-compose -f ${COMPOSE_FILE} up -d --force-recreate client-${i}
    echo "Started client ${i}..."
done

sleep ${SLEEP_INFRA_DURATION}
echo "=== Starting experiments! ==="

for ((i=${MIN_CLIENT_NUM}+1; i<=${MAX_CLIENT_NUM}; i++))
do
    docker-compose -f ${COMPOSE_FILE} up -d --force-recreate client-${i}
    echo "Started client ${i}..."
    sleep ${SLEEP_INFRA_DURATION}
    echo "=== Starting Experiment ==="

    prefix="$OUTPUT_ROOT/${COMPOSE_FILE%.*}-${i}-clients"

    bash bench_test.sh ${ROOT_LB_ADDR} ${prefix} ${IMG_RESIZE_PATH}

    echo "=== Experiment Concluded ==="
    sleep ${SLEEP_TEST_DURATION}
done

echo "=== Done with experiments, killing container ==="
docker-compose -f ${COMPOSE_FILE} down

echo "=== Finished... ==="
