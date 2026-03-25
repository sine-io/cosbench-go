package s3

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestAdapterListObjectsSortsEntries(t *testing.T) {
	var seenPrefix, seenMaxKeys string
	adapter := newTestAdapter(t, "", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Query().Get("list-type") != "2" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
		seenPrefix = r.URL.Query().Get("prefix")
		seenMaxKeys = r.URL.Query().Get("max-keys")
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
  <Contents><Key>pref-b</Key><Size>2</Size></Contents>
  <Contents><Key>pref-a</Key><Size>1</Size></Contents>
</ListBucketResult>`))
	})

	items, err := adapter.ListObjects(context.Background(), "bucket", "pref", 7)
	if err != nil {
		t.Fatal(err)
	}
	if seenPrefix != "pref" || seenMaxKeys != "7" {
		t.Fatalf("prefix=%q maxKeys=%q", seenPrefix, seenMaxKeys)
	}
	if len(items) != 2 || items[0].Key != "pref-a" || items[1].Key != "pref-b" {
		t.Fatalf("items = %#v", items)
	}
}

func TestAdapterMultipartPutAppliesStorageClassToMultipartUpload(t *testing.T) {
	var createMultipartStorageClass string
	partUploads := 0
	adapter := newTestAdapter(t, "STANDARD", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && strings.Contains(r.URL.RawQuery, "uploads"):
			createMultipartStorageClass = r.Header.Get("X-Amz-Storage-Class")
			w.Header().Set("Content-Type", "application/xml")
			_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<InitiateMultipartUploadResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
  <Bucket>bucket</Bucket><Key>key</Key><UploadId>upload-1</UploadId>
</InitiateMultipartUploadResult>`))
		case r.Method == http.MethodPut && strings.Contains(r.URL.RawQuery, "partNumber="):
			_, _ = io.Copy(io.Discard, r.Body)
			partUploads++
			w.Header().Set("ETag", `"etag-1"`)
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodPost && strings.Contains(r.URL.RawQuery, "uploadId="):
			_, _ = io.Copy(io.Discard, r.Body)
			w.Header().Set("Content-Type", "application/xml")
			_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CompleteMultipartUploadResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
  <Bucket>bucket</Bucket><Key>key</Key><ETag>"etag-1"</ETag>
</CompleteMultipartUploadResult>`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	})

	size := int64(5*1024*1024 + 1)
	if err := adapter.MultipartPut(context.Background(), "bucket", "key", bytes.NewReader(make([]byte, size)), size, 5*1024*1024); err != nil {
		t.Fatal(err)
	}
	if createMultipartStorageClass != "STANDARD" {
		t.Fatalf("storage class header = %q", createMultipartStorageClass)
	}
	if partUploads == 0 {
		t.Fatal("expected multipart upload requests")
	}
}

func newTestAdapter(t *testing.T, storageClass string, handler func(http.ResponseWriter, *http.Request)) *Adapter {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(server.Close)

	client := awss3.NewFromConfig(aws.Config{
		Region:      "us-east-1",
		Credentials: credentials.NewStaticCredentialsProvider("ak", "sk", ""),
		HTTPClient:  server.Client(),
	}, func(o *awss3.Options) {
		o.BaseEndpoint = aws.String(server.URL)
		o.UsePathStyle = true
		if storageClass != "" {
			o.APIOptions = append(o.APIOptions, setStorageClassMiddleware(storageClass))
		}
	})

	return &Adapter{
		backend: "s3",
		config:  Config{StorageClass: storageClass},
		client:  client,
	}
}
