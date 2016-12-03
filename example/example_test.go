package example

import (
	"bufio"
	"bytes"
	"github.com/espebra/blobstore"
	"io"
	"os"
	"testing"
)

var (
	key      = "accesskey"
	secret   = "secretkey"
	endpoint = "127.0.0.1:9000"
	location = "us-east-1"
	bucket   = "foobar"
	useSSL   = false
	name     = "bar"
	data     = []byte("some more content")
)

func TestFileSystemProvider(t *testing.T) {
	// Initialize a file system provider
	f := blobstore.NewFileSystemProvider(&blobstore.ProviderData{})

	// Set basedir
	f.Configure(os.TempDir())

	// Write data
	r := io.Reader(
		bytes.NewReader(data),
	)
	bytes, err := f.Store(name, r)
	if err != nil {
		t.Fatalf("Unable to write data: %s\n", err)
	}
	t.Logf("Wrote %d bytes to the file %s in the system provider.", bytes, name)
}

func TestS3Provider(t *testing.T) {
	// Initialize a file system provider
	f := blobstore.NewFileSystemProvider(&blobstore.ProviderData{})

	// Set basedir
	f.Configure(os.TempDir())

	// Copy the file from the file system provider to the S3 provider.
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		bytes, err := f.Retrieve(name, w)
		if err != nil {
			t.Fatalf("Unable to read data: %s\n", err)
		}
		t.Logf("Retrieved %d bytes from the file system provider.", bytes)
	}()

	s := blobstore.NewS3Provider(&blobstore.ProviderData{})
	s.Configure(endpoint, location, bucket, useSSL)
	s.Credentials(key, secret)

	reader := bufio.NewReader(r)
	bytes, err := s.Store(name, reader)
	if err != nil {
		t.Fatalf("Unable to write data: %s\n", err)
	}
	t.Logf("Wrote %d bytes to the file %s in the s3 provider.", bytes, name)

	// Verify that the file exists in S3
	exists, err := s.Exists(name)
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	if exists == false {
		t.Fatal("File does not exist, when it should")
	}

	// Remove the file from the file system
	if err := f.Remove(name); err != nil {
		t.Fatalf("Unable to remove from the file system: %s\n", err)
	}

	// Verify that the file has been removed from the file system
	exists, err = f.Exists(name)
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	if exists != false {
		t.Fatal("File exists, when it should not")
	}

	// Remove the file from S3
	if err := s.Remove(name); err != nil {
		t.Fatalf("Unable to remove from S3: %s\n", err)
	}
}
