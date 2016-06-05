#!/bin/bash

PRIVATE_IP=`ip addr show eth0 | grep inet | head -1 | awk '{print $2}' | cut -f1 -d'/'`

echo "Start ivr-server"

if [ -z "$LOG_LEVEL" ]; then
    LOG_LEVEL=INFO
fi

./ivrtest-server \
    -h "$PRIVATE_IP" \
    -p "$PORT" \
    -n "$PHONE_NUMBER" \
    -r "$RESTCOMM_HOST:$RESTCOMM_PORT" \
    -r-user "$RESTCOMM_USER" \
    -r-pswd "$RESTCOMM_PSWD" \
    -res '' \
    -res-msg $RES_MSG \
    -res-confirm $RES_CONFIRM \
    -l $LOG_LEVEL

