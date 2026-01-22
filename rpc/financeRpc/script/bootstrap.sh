#!/bin/bash
CURDIR=$(cd $(dirname $0); pwd)
BinaryName=financeRpc
echo "$CURDIR/bin/${BinaryName}"
exec $CURDIR/bin/${BinaryName}
