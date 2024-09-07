#!/bin/bash

PORT=9090
RESPONSE=resp.txt

trap "echo Bye!; exit;" SIGINT SIGTERM

while :
do
  echo ""
  echo "-----------------------------------------------------"
  echo ""
  cat $RESPONSE | nc -nlvp $PORT;
done
