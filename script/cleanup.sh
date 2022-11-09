#!/bin/bash

option=$1

cnt=`find ../log -name "*.log" | wc -l`
list=`find ../log -name "*.log"`
echo $cnt
echo "list : $list"

if [ !-z $option ]; then
    echo "remove whole directory"
    rm -rf ../log
else
    rm -rf ../log/*.log
fi

echo "CLEANUP $cnt LOGS"