#!/bin/bash
terraform apply -auto-approve
terraform output -json > output.json
./create_update_alias.sh