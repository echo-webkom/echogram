package images

import (
	"os"
	"strings"
)

func getCredentials() (string, string, string) {
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	containerName := os.Getenv("AZURE_STORAGE_CONTAINER_NAME")

	return accountName, accountKey, containerName
}

func validImageType(filename string) bool {
	validTypes := []string{".jpg", ".jpeg", ".png", ".gif"}
	for _, validType := range validTypes {
		if strings.HasSuffix(filename, validType) {
			return true
		}
	}
	return false
}
