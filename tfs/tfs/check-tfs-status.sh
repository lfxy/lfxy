#!/bin/bash

# Author owner: chenqiangzhishen@163.com

#NS_PORT=8668
#NS_NAME=tfsnameserver

echo $CHECKED_IP $CHECKED_NAME
for check_times in `seq 1 3`; do
    sudo netstat -anp | grep $NS_PORT | grep LISTEN > /dev/null 2>&1
    if [[ $? -eq 1 ]]; then
        echo "111111"
        sleep 1
        if [ $check_times -eq 3 ]; then
            echo "aaaaaa"
            #sudo pkill $NS_NAME
        fi
    fi
done

