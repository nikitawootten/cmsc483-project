#!/usr/bin/env bash
# This file performs load testing on a variable number of clients.
# Given a docker compose file, it will start all load balancers and progressively test the system as more CPUs are added
#
# Ensure that all clients have the form `client-{#}` and all load balancers contain the suffix `-lb` for this script to
# properly find all the relevant clients
# Also sure that the root load balancer is accessible from localhost and runs off port 8080

# The compose file to be tested
COMPOSE_FILE=$1
# The point at which testing begins (used to prevent starving a load balancer)
MIN_CLIENT_NUM=$2
# Number of clients this experiment will test to
MAX_CLIENT_NUM=$3

LOAD_BALANCER_SUFFIX="lb"

SLEEP_INFRA_DURATION=5
# amount of time to sleep between tests
SLEEP_TEST_DURATION=30

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

for (( i=${MIN_CLIENT_NUM}+1; i<=${MAX_CLIENT_NUM}; i++ ))
do
    docker-compose -f ${COMPOSE_FILE} up -d --force-recreate client-${i}
    echo "Started client ${i}..."
    sleep ${SLEEP_INFRA_DURATION}
    echo "=== Starting Experiment ==="
    # todo put experiment here
    echo "=== Experiment Concluded ==="
    sleep ${SLEEP_TEST_DURATION}
done

echo "=== Done with experiments, killing container ==="
docker-compose -f ${COMPOSE_FILE} down

echo "=== Finished... ==="