#!/bin/bash

echo -n "Digits=7" > /tmp/req_gather.data

ab -n 100000 -c 1000 -p /tmp/req_gather.data -T 'application/x-www-form-urlencoded' 127.0.0.1:7090/gather 
