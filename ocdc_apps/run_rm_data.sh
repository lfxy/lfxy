#!/bin/bash

for loop in {1..2000}
do
    echo "running... $loop"
    ./ocdc_apps
    sleep 3
done
