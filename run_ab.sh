#!/bin/bash

echo -n "Digits=7" > /tmp/req_gather.data

curl 127.0.0.1:7090/start
ab -n 100000 -c 500 127.0.0.1:7090/incoming &>/dev/null &
ab -n 100000 -c 500 -p /tmp/req_gather.data -T 'application/x-www-form-urlencoded' 127.0.0.1:7090/gather

curl 127.0.0.1:7090
