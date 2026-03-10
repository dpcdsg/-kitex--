#!/usr/bin/env bash
RUN_NAME="product"
mkdir -p output/bin output/conf
cp script/bootstrap.sh output
cp -r conf/* output/conf 2>/dev/null; true
chmod +x output/bootstrap.sh
go build -o output/bin/${RUN_NAME}
