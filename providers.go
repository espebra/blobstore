package blobstore

import (
	"github.com/espebra/blobstore/common"
	"github.com/espebra/blobstore/filesystem"
	"github.com/espebra/blobstore/s3"
)

// New is used to create and initialize a new storage provider.
func New(provider string, p *common.ProviderData) common.Provider {
	switch provider {
	case "s3":
		return s3.New(p)
	default:
		return filesystem.New(p)
	}
}
