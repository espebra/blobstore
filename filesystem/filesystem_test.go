package filesystem

import (
	"bytes"
	"io"
	"os"
	"testing"
	"github.com/espebra/blobstore"
)

var (
	name = "foo"
	data = []byte("some content")
)

func newFileSystemProvider() *FileSystemProvider {
	return New(&blobstore.ProviderData{})
}

func TestFileSystemProviderDefaults(t *testing.T) {
	dir := os.TempDir()
	p := newFileSystemProvider()
	p.Configure(dir)

	// Verify default values
	if p.BaseDir != dir {
		t.Fatal("Unexpected basedir: ", dir)
	}
}

// Verify that we can store a file
func TestFileSystemProviderStore(t *testing.T) {
	p := newFileSystemProvider()
	p.Configure(os.TempDir())

	r := io.Reader(
		bytes.NewReader(data),
	)
	bytes, err := p.Store(name, r)
	if err != nil {
		t.Fatalf("Unable to write data: %s", err)
	}
	if bytes != 12 {
		t.Fatalf("Unexpected number of bytes stored: %d", bytes)
	}
}

func TestFileSystemProviderExists(t *testing.T) {
	p := newFileSystemProvider()
	p.Configure(os.TempDir())

	exists, err := p.Exists(name)
	if err != nil {
		t.Fatal("Unexpected error: " + err.Error())
	}
	if exists == false {
		t.Fatal("File does not exist, when it should")
	}
}

// Verify that we can read a file
func TestFileSystemProviderRetrieve(t *testing.T) {
	p := newFileSystemProvider()
	p.Configure(os.TempDir())

	var buf bytes.Buffer
	bytes, err := p.Retrieve(name, &buf)
	if err != nil {
		t.Fatal("Unable to read data.")
	}
	if bytes != 12 {
		t.Fatalf("Unexpected number of bytes retrieved: %d", bytes)
	}
	if buf.String() != string(data) {
		t.Fatalf("Unexpected contents: %s\n", buf.String())
	}
}

// Verify that we can remove a file
func TestFileSystemProviderRemove(t *testing.T) {
	p := newFileSystemProvider()
	p.Configure(os.TempDir())

	err := p.Remove(name)
	if err != nil {
		t.Fatal("Unable to remove data.")
	}

	// Try to read the data that was just removed. It should fail.
	var buf bytes.Buffer
	bytes, err := p.Retrieve(name, &buf)
	if err == nil {
		t.Fatal("Unexpected success")
	}
	if bytes != 0 {
		t.Fatalf("Unexpected number of bytes retrieved: %d", bytes)
	}

	exists, err := p.Exists(name)
	if err != nil {
		t.Fatal("Unexpected error: " + err.Error())
	}
	if exists == true {
		t.Fatal("File exists, when it should not")
	}
	if buf.String() != "" {
		t.Fatalf("Unexpected contents: %s\n", buf.String())
	}
}
