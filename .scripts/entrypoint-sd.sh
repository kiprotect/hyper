#!/bin/sh
echo "Execute sd with user hyper"
exec su hyper -c "./sd $*"
