package blobstore

import (
	"bufio"
	"io"
	"os"
	"path"
)

type FileSystemProvider struct {
	*ProviderData
	BaseDir string
}

func NewFileSystemProvider(p *ProviderData) *FileSystemProvider {
	p.Encryption = false
	p.Secret = ""

	return &FileSystemProvider{
		ProviderData: p,
		BaseDir:      "/srv/blobstore",
	}
}

// Store named file
func (p *FileSystemProvider) Store(name string, data io.Reader) (int64, error) {
	fpath := path.Join(p.BaseDir, name)
	f, err := os.Create(fpath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	b := bufio.NewReader(data)
	bytes, err := b.WriteTo(f)
	return bytes, err
}

// Retrieve named file
func (p *FileSystemProvider) Retrieve(name string, fp io.Writer) (int64, error) {
	fpath := path.Join(p.BaseDir, name)
	f, err := os.Open(fpath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	bytes, err := io.Copy(fp, f)
	return bytes, err
}

// Delete named file
func (p *FileSystemProvider) Delete(name string) error {
	fpath := path.Join(p.BaseDir, name)
	return os.Remove(fpath)
}

// Does a named file exist
func (p *FileSystemProvider) Exists(name string) (bool, error) {
	fpath := path.Join(p.BaseDir, name)
	_, err := os.Stat(fpath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
