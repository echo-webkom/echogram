package services

type BlobManager interface {
	// Get retrieves a blob from the blob storage
	Get(filename string) ([]byte, error)
	// Add adds a blob to the blob storage, should not overwrite existing blobs
	Add(filename string, file []byte) error
}
