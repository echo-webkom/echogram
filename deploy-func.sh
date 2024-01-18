#!/bin/bash

FUNCTION_APP_NAME=torger-function-app
RESOURCE_GROUP_NAME=torger-resources

if [ -f handler ]; then
    rm main
fi

# build go binary / request handler
GOOS=linux GOARCH=amd64 go build -o main

# deploy function app
func azure functionapp publish $FUNCTION_APP_NAME --custom

# set environment variables
if [ ! -f ./.env ]; then
    echo "File .env does not exist"
    exit 1
fi

echo "Setting environment variables..."

source ./.env
AZURE_STORAGE_ACCOUNT_NAME=$AZURE_STORAGE_ACCOUNT_NAME
AZURE_STORAGE_ACCOUNT_KEY=$AZURE_STORAGE_ACCOUNT_KEY
AZURE_STORAGE_CONTAINER_NAME=$AZURE_STORAGE_CONTAINER_NAME

if [ -z "$AZURE_STORAGE_ACCOUNT_NAME" ]; then
    echo "AZURE_STORAGE_ACCOUNT_NAME is not set"
    exit 1
fi

if [ -z "$AZURE_STORAGE_ACCOUNT_KEY" ]; then
    echo "AZURE_STORAGE_ACCOUNT_KEY is not set"
    exit 1
fi

if [ -z "$AZURE_STORAGE_CONTAINER_NAME" ]; then
    echo "AZURE_STORAGE_CONTAINER_NAME is not set"
    exit 1
fi

az functionapp config appsettings set \
    --name $FUNCTION_APP_NAME \
    --resource-group $RESOURCE_GROUP_NAME \
    --settings \
        AZURE_STORAGE_ACCOUNT_NAME=$AZURE_STORAGE_ACCOUNT_NAME \
        AZURE_STORAGE_ACCOUNT_KEY=$AZURE_STORAGE_ACCOUNT_KEY \
        AZURE_STORAGE_CONTAINER_NAME=$AZURE_STORAGE_CONTAINER_NAME \
    --output none

echo "Success"