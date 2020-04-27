#!/usr/bin/env bash
# get docker container's IP address!
IP_ADDR=$(awk 'END{print $1}' /etc/hosts) /dist/main "$@"