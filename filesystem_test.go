package blobstore

import (
	"bytes"
	"io"
	"testing"
)

func newFileSystemProvider() *FileSystemProvider {
	return NewFileSystemProvider(
		&ProviderData{})
}

func TestFileSystemProviderDefaults(t *testing.T) {
	p := newFileSystemProvider()

	// Verify default values
	if p.BaseDir != "/srv/blobstore" {
		t.Fatal("Unexpected basedir: ", p.BaseDir)
	}
}

// Verify that we can store a file
func TestFileSystemProviderStore(t *testing.T) {
	p := newFileSystemProvider()
	r := io.Reader(
		bytes.NewReader([]byte("some content")),
	)
	bytes, err := p.Store("foo", r)
	if err != nil {
		t.Fatal("Unable to write data.")
	}
	if bytes != 12 {
		t.Fatalf("Unexpected number of bytes stored: %d", bytes)
	}

	exists, err := p.Exists("foo")
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

	var f bytes.Buffer
	bytes, err := p.Retrieve("foo", &f)
	if err != nil {
		t.Fatal("Unable to read data.")
	}
	if bytes != 12 {
		t.Fatalf("Unexpected number of bytes retrieved: %d", bytes)
	}
}

// Verify that we can delete a file
func TestFileSystemProviderDelete(t *testing.T) {
	p := newFileSystemProvider()

	err := p.Delete("foo")
	if err != nil {
		t.Fatal("Unable to delete data.")
	}

	// Try to read the data that was just deleted. It should fail.
	var f bytes.Buffer
	bytes, err := p.Retrieve("foo", &f)
	if err == nil {
		t.Fatal("Unexpected success")
	}
	if bytes != 0 {
		t.Fatalf("Unexpected number of bytes retrieved: %d", bytes)
	}

	exists, err := p.Exists("foo")
	if err != nil {
		t.Fatal("Unexpected error: " + err.Error())
	}
	if exists == true {
		t.Fatal("File exists, when it should not")
	}
}
