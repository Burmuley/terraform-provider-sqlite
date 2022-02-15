#!/usr/bin/env bash

rm -f .terraform.lock.hcl
terraform init
terraform apply

