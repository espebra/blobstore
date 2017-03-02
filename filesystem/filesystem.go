package filesystem

import (
	"bufio"
	"github.com/espebra/blobstore/common"
	"io"
	"os"
	"path"
	"path/filepath"
	"errors"
)

type FileSystemProvider struct {
	*common.ProviderData
	baseDir string
}

// NewFileSystemProvider initializes a new FileSystemProvider with default
// values.
func New(p *common.ProviderData) *FileSystemProvider {
	p.Encryption = false
	p.Secret = ""

	return &FileSystemProvider{ProviderData: p}
}

// Setup
func (p *FileSystemProvider) Setup(cfg map[string]string) error {
	p.baseDir = cfg["basedir"]
	if p.baseDir == "" {
		p.baseDir = "/var/lob/blobstore"
	}
	return nil
}

// Store named file in FileSystemProvider. The return value bytes is the number
// of bytes that was stored.
func (p *FileSystemProvider) Store(name string, data io.Reader) (bytes int64, err error) {
	fpath := path.Join(p.baseDir, name)
	f, err := os.Create(fpath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	b := bufio.NewReader(data)
	bytes, err = b.WriteTo(f)
	return bytes, err
}

// Retrieve named file from FileSystemProvider. The return value bytes is the
// number of bytes that was retrieved.
func (p *FileSystemProvider) Retrieve(name string, fp io.Writer) (bytes int64, err error) {
	if name != filepath.Base(name) {
		return 0, errors.New("Invalid name")
	}

	fpath := path.Join(p.baseDir, name)
	f, err := os.Open(fpath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	bytes, err = io.Copy(fp, f)
	return bytes, err
}

// Remove named file from FileSystemProvider.
func (p *FileSystemProvider) Remove(name string) error {
	fpath := path.Join(p.baseDir, name)
	return os.Remove(fpath)
}

// Exists will verify if a named file exists in FileSystemProvider. The return
// value exists is a boolean indicating if the named file exists or not.
func (p *FileSystemProvider) Exists(name string) (exists bool, err error) {
	fpath := path.Join(p.baseDir, name)
	_, err = os.Stat(fpath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
