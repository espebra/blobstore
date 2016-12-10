package s3

import (
	"bytes"
	"github.com/espebra/blobstore/common"
	"io"
	"testing"
)

var (
	key      = "accesskey"
	secret   = "secretkey"
	endpoint = "127.0.0.1:9000"
	location = "us-east-1"
	bucket   = "foobar"
	useSSL   = "no"
	cfg      = map[string]string{}
)

func new() (*S3Provider, error) {
	cfg := map[string]string{}
	cfg["key"] = key
	cfg["secret"] = secret
	cfg["endpoint"] = endpoint
	cfg["location"] = location
	cfg["bucket"] = bucket
	cfg["useSSL"] = useSSL

	p := New(&common.ProviderData{})
	err := p.Setup(cfg)
	return p, err

}

func TestS3ProviderDefaults(t *testing.T) {
	p, err := new()
	if err != nil {
		t.Fatal("Unable to create provider: ", err.Error())
	}
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
	p, err := new()
	if err != nil {
		t.Fatal("Unable to create provider: ", err.Error())
	}

	r := io.Reader(
		bytes.NewReader([]byte("some content")),
	)
	nBytes, err := p.Store("foo", r)
	if err != nil {
		t.Fatalf("Unable to write data: %s\n", err)
	}
	if nBytes != 12 {
		t.Fatalf("Unexpected number of bytes stored: %d\n", nBytes)
	}
}

func TestS3ProviderExists(t *testing.T) {
	p, err := new()
	if err != nil {
		t.Fatal("Unable to create provider: ", err.Error())
	}

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
	p, err := new()
	if err != nil {
		t.Fatal("Unable to create provider: ", err.Error())
	}

	var buf bytes.Buffer
	nBytes, err := p.Retrieve("foo", &buf)
	if err != nil {
		t.Fatalf("Unable to read data: %s\n", err)
	}
	if nBytes != 12 {
		t.Fatalf("Unexpected number of bytes retrieved: %d\n", nBytes)
	}
	if buf.String() != "some content" {
		t.Fatalf("Unexpected content: %s\n", buf.String())
	}
}

// Verify that we can remove a file
func TestS3ProviderRemove(t *testing.T) {
	p, err := new()
	if err != nil {
		t.Fatal("Unable to create provider: ", err.Error())
	}

	if err := p.Remove("foo"); err != nil {
		t.Fatalf("Unable to remove data: %s\n", err)
	}

	// Try to read the data that was just removed. It should fail.
	var buf bytes.Buffer
	nBytes, err := p.Retrieve("foo", &buf)
	if err == nil {
		t.Fatal("Unexpected success")
	}
	if nBytes != 0 {
		t.Fatalf("Unexpected number of bytes retrieved: %d\n", nBytes)
	}

	exists, err := p.Exists("foo")
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	if exists == true {
		t.Fatal("File exists, when it should not")
	}
	if buf.String() == "some content" {
		t.Fatalf("Unexpected content: %s\n", buf.String())
	}
}
