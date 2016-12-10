package common

import (
	"io"
)

type Provider interface {
	Data() *ProviderData
	Store(string, io.Reader) (int64, error)
	Retrieve(string, io.Writer) (int64, error)
	Remove(string) error
	Exists(string) (bool, error)
	Setup(map[string]string) error
}

type ProviderData struct {
	Encryption bool
	Secret     string
}

func (p *ProviderData) Data() *ProviderData {
	return p
}
