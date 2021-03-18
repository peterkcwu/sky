#!/bin/sh

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

PIDS=$(GetRunningPID)

if [[ -n "$PIDS" ]];then
    ps -fl $PIDS
else
    echo "no running process"
    exit 1
fi