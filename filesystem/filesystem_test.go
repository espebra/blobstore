package filesystem

import (
	"bytes"
	"io"
	"os"
	"testing"
	"github.com/espebra/blobstore/common"
)

var (
	dir = os.TempDir()
)

func new() (*FileSystemProvider, error) {
        cfg := map[string]string{}
        cfg["basedir"] = dir

        p := New(&common.ProviderData{})
        err := p.Setup(cfg)
        return p, err
}

func TestFileSystemProviderDefaults(t *testing.T) {
        p, err := new()
        if err != nil {
                t.Fatal("Unable to create provider: ", err.Error())
        }

	// Verify default values
	if p.baseDir != dir {
		t.Fatal("Unexpected basedir: ", dir)
	}
}

// Verify that we can store a file
func TestFileSystemProviderStore(t *testing.T) {
        p, err := new()
        if err != nil {
                t.Fatal("Unable to create provider: ", err.Error())
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
}

func TestFileSystemProviderExists(t *testing.T) {
        p, err := new()
        if err != nil {
                t.Fatal("Unable to create provider: ", err.Error())
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
        p, err := new()
        if err != nil {
                t.Fatal("Unable to create provider: ", err.Error())
        }

	var buf bytes.Buffer
	nBytes, err := p.Retrieve("foo", &buf)
	if err != nil {
		t.Fatal("Unable to read data.")
	}
	if nBytes != 12 {
		t.Fatalf("Unexpected number of bytes retrieved: %d", nBytes)
	}
	if buf.String() != "some content" {
		t.Fatalf("Unexpected content: %s\n", buf.String())
	}
}

// Verify that we can remove a file
func TestFileSystemProviderRemove(t *testing.T) {
        p, err := new()
        if err != nil {
                t.Fatal("Unable to create provider: ", err.Error())
        }

	if err := p.Remove("foo"); err != nil {
		t.Fatal("Unable to remove data.")
	}

	// Try to read the data that was just removed. It should fail.
	var buf bytes.Buffer
	nBytes, err := p.Retrieve("foo", &buf)
	if err == nil {
		t.Fatal("Unexpected success")
	}
	if nBytes != 0 {
		t.Fatalf("Unexpected number of bytes retrieved: %d", nBytes)
	}

	exists, err := p.Exists("foo")
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
