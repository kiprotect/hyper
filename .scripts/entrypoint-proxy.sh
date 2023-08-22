#!/bin/sh
echo "Execute proxy with user hyper"
exec su hyper -c "./proxy $*"
