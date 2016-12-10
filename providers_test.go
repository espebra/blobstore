package blobstore

import (
	"bytes"
	"github.com/espebra/blobstore/common"
	"io"
	"os"
	"testing"
)

func TestNewProvider(t *testing.T) {
	p := New("filesystem", &common.ProviderData{})
	cfg := map[string]string{}
	cfg["basedir"] = os.TempDir()
	if err := p.Setup(cfg); err != nil {
		t.Fatal("Setup failed: ", err.Error())
	}

	r := io.Reader(
		bytes.NewReader([]byte("some content")),
	)
	nBytes, err := p.Store("foo", r)
	if err != nil {
		t.Fatal("Unable to write data.")
	}
	if nBytes != 12 {
		t.Fatalf("Unexpected number of bytes stored: %d", nBytes)
	}

	exists, err := p.Exists("foo")
	if err != nil {
		t.Fatal("Unexpected error: " + err.Error())
	}
	if exists == false {
		t.Fatal("File does not exist, when it should")
	}

	if err := p.Remove("foo"); err != nil {
		t.Fatal("Unable to remove data.")
	}

	// Try to read the data that was just removed. It should fail.
	var buf bytes.Buffer
	nBytes, err = p.Retrieve("foo", &buf)
	if err == nil {
		t.Fatal("Unexpected success")
	}
	if nBytes != 0 {
		t.Fatalf("Unexpected number of bytes retrieved: %d", nBytes)
	}

	exists, err = p.Exists("foo")
	if err != nil {
		t.Fatal("Unexpected error: " + err.Error())
	}
	if exists == true {
		t.Fatal("File exists, when it should not")
	}
	if buf.String() == "some content" {
		t.Fatalf("Unexpected content: %s\n", buf.String())
	}
}
