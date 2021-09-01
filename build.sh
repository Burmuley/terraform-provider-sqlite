#!/usr/bin/env bash

rm -f terraform-provider-sqlite
go build -o terraform-provider-sqlite
mkdir -p ~/.terraform.d/plugins/burmuley.com/edu/sqlite/0.1/darwin_amd64
mv terraform-provider-sqlite ~/.terraform.d/plugins/burmuley.com/edu/sqlite/0.1/darwin_amd64
