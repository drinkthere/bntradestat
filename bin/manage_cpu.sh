#!/bin/bash

cpuidx=8

start() {
    echo "Starting the process @cpu '$cpuidx'..."
    nohup taskset -c "$cpuidx"  ./bntradestat ../config/config.json >> /data/will/bntradestat/nohup.log 2>&1 &
    echo "Process started."
}

stop() {
    echo "Stopping the process for account '$account'..."
    pid=$(pgrep -f "bntradestat ../config/config.json")
    if [ -n "$pid" ]; then
        kill -SIGINT $pid
        echo "Process stopping: "
        sleep 10
        echo "Process stopped."
    else
        echo "Process is not running."
    fi
}

restart() {
    stop
    sleep 10
    start
}

case "$2" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    *)
        echo "Usage: $0 account {start|stop|restart}"
        exit 1
        ;;
esac

exit 0
