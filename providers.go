package blobstore

import (
	"io"
)

type Provider interface {
	Data() *ProviderData
	Store(string, io.Reader) (int64, error)
	Retrieve(string, io.Writer) (int64, error)
	Delete(string) error
	Exists(string) (bool, error)
}

func New(provider string, p *ProviderData) Provider {
	switch provider {
	case "s3":
		return NewS3Provider(p)
	default:
		return NewFileSystemProvider(p)
	}
}
