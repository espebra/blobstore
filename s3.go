package blobstore

import (
	"bufio"
	"errors"
	"github.com/minio/minio-go"
	"io"
	"log"
)

type S3Provider struct {
	*ProviderData
	key      string
	secret   string
	useSSL   bool
	endpoint string
	location string
	bucket   string
}

func NewS3Provider(p *ProviderData) *S3Provider {
	p.Encryption = false
	p.Secret = ""

	return &S3Provider{ProviderData: p}
}

func (p *S3Provider) Credentials(key, secret string) {
	p.key = key
	p.secret = secret
}

func (p *S3Provider) Configure(endpoint, location, bucket string, useSSL bool) {
	p.endpoint = endpoint
	if p.endpoint == "" {
		p.endpoint = "s3.amazonaws.com"
	}
	p.location = location
	if p.location == "" {
		p.location = "us-east-1"
	}
	p.bucket = bucket
	if p.bucket == "" {
		p.bucket = "blobstore"
	}
	p.useSSL = useSSL
}

// Store named file
func (p *S3Provider) Store(name string, data io.Reader) (int64, error) {
	minioClient, err := minio.New(p.endpoint, p.key, p.secret, p.useSSL)
	if err != nil {
		return 0, err
	}
	err = minioClient.MakeBucket(p.bucket, p.location)
	if err == nil {
		log.Printf("Successfully created bucket %s\n", p.bucket)
	} else {
		exists, err := minioClient.BucketExists(p.bucket)
		if err != nil {
			return 0, err
		}
		if exists == false {
			return 0, errors.New("Unable to create bucket")
		}
	}

	contentType := "application/octet-stream"
	b := bufio.NewReader(data)
	bytes, err := minioClient.PutObject(p.bucket, name, b, contentType)
	if err != nil {
		return 0, err
	}
	return bytes, err
}

// Retrieve named file
func (p *S3Provider) Retrieve(name string, fp io.Writer) (int64, error) {
	minioClient, err := minio.New(p.endpoint, p.key, p.secret, p.useSSL)
	if err != nil {
		log.Fatalf("Unable to create client: %s\n", err)
		return 0, errors.New("Unable to create client")
	}
	f, err := minioClient.GetObject(p.bucket, name)
	if err != nil {
		log.Fatalf("Unabel to get %s: %s\n", name, err)
		return 0, errors.New("Unable to retrieve " + name + ":" + err.Error())
	}
	bytes, err := io.Copy(fp, f)
	return bytes, err
}

// Delete named file
func (p *S3Provider) Delete(name string) error {
	minioClient, err := minio.New(p.endpoint, p.key, p.secret, p.useSSL)
	if err != nil {
		return err
	}
	return minioClient.RemoveObject(p.bucket, name)
}

// Does a named file exist
func (p *S3Provider) Exists(name string) (bool, error) {
	minioClient, err := minio.New(p.endpoint, p.key, p.secret, p.useSSL)
	if err != nil {
		return false, err
	}
	_, err = minioClient.StatObject(p.bucket, name)
	if err == nil {
		return true, nil
	}
	return false, nil
}
