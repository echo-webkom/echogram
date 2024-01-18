#!/bin/bash

OUT_FILE="plan.out"

set -e

terraform init && \
terraform fmt -recursive && \
terraform validate

if [ -f "$OUT_FILE" ]; then
  rm $OUT_FILE
fi

echo "Creating $OUT_FILE file..."
terraform plan -out $OUT_FILE

echo "Applying $OUT_FILE file..."
terraform apply "$OUT_FILE"

rm $OUT_FILE