package blobstore

type ProviderData struct {
	Encryption bool
	Secret     string
}

func (p *ProviderData) Data() *ProviderData { return p }
