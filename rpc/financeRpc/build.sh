#!/bin/bash
RUN_NAME="financeRpc"
mkdir -p output/bin
cp script/bootstrap.sh output 2>/dev/null || :
chmod +x output/bootstrap.sh
go build -o output/bin/${RUN_NAME}
