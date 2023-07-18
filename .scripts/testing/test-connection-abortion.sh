#!/bin/bash
while true
do
    HYPER_SETTINGS=settings/dev/roles/private-proxy-hyper-1 hyper --level debug server run &
    RUNNING_PID=$!
    sleep 0.1
    kill ${RUNNING_PID}
done
