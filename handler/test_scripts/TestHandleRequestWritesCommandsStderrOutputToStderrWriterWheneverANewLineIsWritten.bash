#!/usr/bin/env bash

set -o errtrace

echo "$$" >&2

echo -n "stdout"


loop_forever() {
  while true; do
    sleep 1
  done
}



loop_forever &
child="$!"

counter=0
log_line() {
  echo "log line $counter"
  (( counter++ ))
}

sigterm_trap() {
  kill -TERM "$child"
  exit 0
}

trap 'sigterm_trap' SIGTERM
trap 'log_line' SIGUSR1

while true; do
  wait "$child"
done
