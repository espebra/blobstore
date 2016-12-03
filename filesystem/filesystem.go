package filesystem

import (
	"bufio"
	"io"
	"os"
	"path"
	"github.com/espebra/blobstore"
)

type FileSystemProvider struct {
	*blobstore.ProviderData
	BaseDir string
}

// New initializes a new FileSystemProvider with default values.
func New(p *blobstore.ProviderData) *FileSystemProvider {
	p.Encryption = false
	p.Secret = ""

	return &FileSystemProvider{ProviderData: p}
}

// Configure configures a FileSystemProvider.
func (p *FileSystemProvider) Configure(basedir string) {
	p.BaseDir = basedir
	if p.BaseDir == "" {
		p.BaseDir = "/srv/blobstore"
	}
}

// Store named file in FileSystemProvider. The return value bytes is the number
// of bytes that was stored.
func (p *FileSystemProvider) Store(name string, data io.Reader) (bytes int64, err error) {
	fpath := path.Join(p.BaseDir, name)
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
	fpath := path.Join(p.BaseDir, name)
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
	fpath := path.Join(p.BaseDir, name)
	return os.Remove(fpath)
}

// Exists will verify if a named file exists in FileSystemProvider. The return
// value exists is a boolean indicating if the named file exists or not.
func (p *FileSystemProvider) Exists(name string) (exists bool, err error) {
	fpath := path.Join(p.BaseDir, name)
	_, err = os.Stat(fpath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
