#!/bin/sh

ARGS=
if [ -n "$LISTEN" ]; then
    ARGS="$ARGS -listen $LISTEN"
fi

/net-delay-time-exporter $ARGS $(echo $SERVERS | sed 's/,/ /g')