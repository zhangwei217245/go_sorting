#!/bin/bash

echo "sleep 20s"

sleep 20s

echo "creating /tmp/go_sorting"

mkdir -p /tmp/go_sorting

echo "run consumer now..."

nohup /app/main -func=con -count=100000000 -numchunk=10 > /tmp/go_sorting/con.log 2>&1 &

echo "run producer now..."

nohup /app/main -func=pro -count=100000000 > /tmp/go_sorting/pro.log 2>&1 &

echo "Please see program logs: /tmp/go_sorting/*.log"

