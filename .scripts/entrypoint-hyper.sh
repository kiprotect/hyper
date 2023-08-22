#!/bin/sh
echo "Execute hyper with user hyper"
exec su hyper -c "./hyper $*"
