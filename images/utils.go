package images

import (
	"os"

	"github.com/echo-webkom/echo-blob/services"
)

func getBlobManager() (services.BlobManager, error) {
	if os.Getenv("ENV") == "dev" {
		lm, err := services.NewLocalBlobManager()
		if err != nil {
			return nil, err
		}

		return lm, nil
	}

	am, err := services.NewAzureBlobManager()
	if err != nil {
		return nil, err
	}

	return am, nil
}
