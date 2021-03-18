#!/bin/bash

export PATH=$PATH:/sbin:/usr/sbin:/usr/local/bin:/usr/bin:/bin:/usr/local/sbin

CURFILE=$(readlink -f "$0")

CURDIR=$(dirname $CURFILE)

MODULE_DIR=${CURDIR%"/tools/op"}

MODULE=$(basename $MODULE_DIR)

if [[ ! -f "$MODULE_DIR/bin/$MODULE" ]];then
    MODULE=$(echo $MODULE | sed 's/[0-9]*$//')
    if [[ ! -f "$MODULE_DIR/bin/$MODULE" ]];then
        echo "unsupport proc"
        exit 1
    fi
fi

LOCKFILE=${CURDIR}/.lock

function GetLock() {
    if [ -e ${LOCKFILE} ] && kill -0 `cat ${LOCKFILE}` 2> /dev/null; then
        echo "cannot run mutiple $0 at the same time"
        exit 1
    fi

    trap "rm -f ${LOCKFILE};exit" INT TERM EXIT
    echo $$ > ${LOCKFILE}
}

function GetRunningPID() {
    PIDS=""
    SUSPECT_PIDS=$(ps -fle | grep "\./$MODULE" | grep -v grep | awk '{print $4}')
    for pid in $SUSPECT_PIDS;do
        EXE=$(readlink -f /proc/$pid/exe)
        if [[ "$EXE" == "${MODULE_DIR}/bin/$MODULE" ]] || [[ "$EXE" == "${MODULE_DIR}/bin/$MODULE (deleted)" ]]; then
            if [[ -z "$PIDS" ]];then
                PIDS="$pid"
            else
                PIDS="$PIDS $pid"
            fi
        fi
    done
    echo $PIDS
}

GetLock

RUNNING_PID=$(GetRunningPID)

if [[ -z "$RUNNING_PID" ]];then
    echo "${MODULE} not running!"
    exit 0
fi

kill ${RUNNING_PID}

MAX_RETRY=100

for ((i=1; i<=$MAX_RETRY; i++))
do
    RUNNING_PID=$(GetRunningPID)

    if [[ -z "$RUNNING_PID" ]];then
        exit 0
    fi

    if [[ "$i" -lt $MAX_RETRY ]];then
        echo "${MODULE} still running pid=$RUNNING_PID, recheck 3 seconds later!"
        sleep 3
    else
        echo "stop ${MODULE} failed, pid=$RUNNING_PID"
        exit 1
    fi
done