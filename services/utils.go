package services

import "os"

type AzureCredentials struct {
	AccountName   string
	AccountKey    string
	ContainerName string
}

func getAzureCredentials() AzureCredentials {
	return AzureCredentials{
		AccountName:   os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"),
		AccountKey:    os.Getenv("AZURE_STORAGE_ACCOUNT_KEY"),
		ContainerName: os.Getenv("AZURE_STORAGE_CONTAINER_NAME"),
	}
}
