#!/bin/bash

while [ 1 ]
do
	./httptest -protocol http -host wordpress.jam10000bo.com -method post -port 8099 -path /cloud2team -count 102
	sleep 1
done
