package services

import (
	"errors"
	"io"
	"os"
)

var (
	ErrBlobAlreadyExists  = errors.New("blob already exists")
	ErrFailedToCreateBlob = errors.New("failed to create blob")
	ErrFailedToWriteBlob  = errors.New("failed to write blob")
)

type LocalBlobManager struct {
	Folder string
}

func NewLocalBlobManager() (*LocalBlobManager, error) {
	folderName := "./local-bucket"

	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		err := os.Mkdir(folderName, 0755)
		if err != nil {
			return nil, err
		}
	}

	return &LocalBlobManager{
		Folder: folderName,
	}, nil
}

func (s *LocalBlobManager) Get(filename string) ([]byte, error) {
	file, err := os.Open(s.Folder + "/" + filename)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *LocalBlobManager) Add(filename string, file []byte) error {
	newFilePath := s.Folder + "/" + filename

	if _, err := os.Stat(newFilePath); err == nil {
		return ErrBlobAlreadyExists
	}

	newFile, err := os.Create(newFilePath)
	if err != nil {
		return ErrFailedToCreateBlob
	}

	_, err = newFile.Write(file)
	if err != nil {
		return ErrFailedToWriteBlob
	}

	return nil
}
