#!/usr/bin/env bash

if [[ "$#" -ne 1 ]]; then
   echo "Usage: ./experiment-suite.sh <output dir>"
   exit 2
fi

OUTPUT_ROOT=$1

experiments=(
    "t1-least-connections-1cpu.yaml"
    "t1-least-connections-2cpu.yaml"
    "t1-least-connections-4cpu.yaml"
    "t1-random-1cpu.yaml"
    "t1-random-2cpu.yaml"
    "t1-random-4cpu.yaml"
    "t1-round-robin-1cpu.yaml"
    "t1-round-robin-2cpu.yaml"
    "t1-round-robin-4cpu.yaml"
    "t2-overlap-1cpu.yaml")
min_client_numbers=(0  0 0 0  0 0 0  0 0 1)
max_client_numbers=(11 5 2 11 5 2 11 5 2 8)

mkdir -p ${OUTPUT_ROOT}

for ((i=0; i<${#experiments[*]}; i++)); do
    echo "--- Starting Experiment For ${experiments[$i]%.*} ---"
    bash load-experiment.sh ${experiments[$i]} ${min_client_numbers[$i]} ${max_client_numbers[$i]} ${OUTPUT_ROOT}
done
