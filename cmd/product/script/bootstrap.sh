#!/usr/bin/env bash
CURDIR=$(cd $(dirname $0); pwd)
BinaryName=product
echo "start $BinaryName"
exec $CURDIR/bin/$BinaryName -config $CURDIR/conf
