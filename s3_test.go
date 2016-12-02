package blobstore

import (
	"bytes"
	"io"
	"testing"
)

var (
	key      = "accesskey"
	secret   = "secretkey"
	endpoint = "127.0.0.1:9000"
	location = "us-east-1"
	bucket   = "foobar"
	useSSL   = false
)

func newS3Provider() *S3Provider {
	return NewS3Provider(&ProviderData{})
}

func TestS3ProviderDefaults(t *testing.T) {
	p := newS3Provider()
	p.Credentials(key, secret)
	p.Configure(endpoint, location, bucket, useSSL)

	if p.key != key {
		t.Fatal("Unexpected key: ", p.key)
	}
	if p.secret != secret {
		t.Fatal("Unexpected secret: ", p.secret)
	}
	if p.endpoint != endpoint {
		t.Fatal("Unexpected endpoint: ", p.endpoint)
	}
	if p.bucket != bucket {
		t.Fatal("Unexpected bucket: ", p.bucket)
	}
	if p.location != location {
		t.Fatal("Unexpected location: ", p.location)
	}
}

// Verify that we can store a file
func TestS3ProviderStore(t *testing.T) {
	p := newS3Provider()
	p.Credentials(key, secret)
	p.Configure(endpoint, location, bucket, useSSL)

	r := io.Reader(
		bytes.NewReader([]byte("some content")),
	)
	bytes, err := p.Store("foo", r)
	if err != nil {
		t.Fatalf("Unable to write data: %s\n", err)
	}
	if bytes != 12 {
		t.Fatalf("Unexpected number of bytes stored: %d\n", bytes)
	}
}

func TestS3ProviderExists(t *testing.T) {
	p := newS3Provider()
	p.Credentials(key, secret)
	p.Configure(endpoint, location, bucket, useSSL)

	exists, err := p.Exists("foo")
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	if exists == false {
		t.Fatal("File does not exist, when it should")
	}
}

// Verify that we can read a file
func TestS3ProviderRetrieve(t *testing.T) {
	p := newS3Provider()
	p.Credentials(key, secret)
	p.Configure(endpoint, location, bucket, useSSL)

	var f bytes.Buffer
	bytes, err := p.Retrieve("foo", &f)
	if err != nil {
		t.Fatalf("Unable to read data: %s\n", err)
	}
	if bytes != 12 {
		t.Fatalf("Unexpected number of bytes retrieved: %d\n", bytes)
	}
}

// Verify that we can delete a file
func TestS3ProviderDelete(t *testing.T) {
	p := newS3Provider()
	p.Credentials(key, secret)
	p.Configure(endpoint, location, bucket, useSSL)

	err := p.Delete("foo")
	if err != nil {
		t.Fatalf("Unable to delete data: %s\n", err)
	}

	// Try to read the data that was just deleted. It should fail.
	var f bytes.Buffer
	bytes, err := p.Retrieve("foo", &f)
	if err == nil {
		t.Fatal("Unexpected success")
	}
	if bytes != 0 {
		t.Fatalf("Unexpected number of bytes retrieved: %d\n", bytes)
	}

	exists, err := p.Exists("foo")
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	if exists == true {
		t.Fatal("File exists, when it should not")
	}
}
