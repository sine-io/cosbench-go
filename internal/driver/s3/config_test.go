package s3

import (
	"errors"
	"net/http"
	"testing"
)

func TestParseConfigS3(t *testing.T) {
	cfg, err := ParseConfig("s3", "endpoint=http://localhost:9000;region=us-east-1;accesskey=ak;secretkey=sk")
	if err != nil {
		t.Fatalf("ParseConfig(): %v", err)
	}
	if cfg.Endpoint != "http://localhost:9000" {
		t.Fatalf("endpoint = %q", cfg.Endpoint)
	}
	if cfg.PathStyle {
		t.Fatal("expected path style disabled by default for s3")
	}
}

func TestParseConfigSIOEnablesPathStyle(t *testing.T) {
	cfg, err := ParseConfig("sio", "endpoint=http://localhost:9000;accesskey=ak;secretkey=sk")
	if err != nil {
		t.Fatalf("ParseConfig(): %v", err)
	}
	if !cfg.PathStyle {
		t.Fatal("expected path style enabled for sio")
	}
}

func TestParseConfigS3UsesLegacyDefaultEndpoint(t *testing.T) {
	cfg, err := ParseConfig("s3", "accesskey=ak;secretkey=sk")
	if err != nil {
		t.Fatalf("ParseConfig(): %v", err)
	}
	if cfg.Endpoint != "http://s3.amazonaws.com" {
		t.Fatalf("endpoint = %q", cfg.Endpoint)
	}
}

func TestParseConfigCompatibilityProfilesKeepLegacyPathStyleDefault(t *testing.T) {
	for _, backend := range []string{"siov1", "gdas"} {
		cfg, err := ParseConfig(backend, "endpoint=http://localhost:9000;accesskey=ak;secretkey=sk")
		if err != nil {
			t.Fatalf("%s ParseConfig(): %v", backend, err)
		}
		if cfg.PathStyle {
			t.Fatalf("%s pathStyle = true", backend)
		}
	}
}

func TestParseConfigRecognizesCompatibilityFields(t *testing.T) {
	cfg, err := ParseConfig("sio", "endpoint=http://localhost:9000;aws_region=cn-east-1;accesskey=ak;secretkey=sk;proxyhost=proxy.local;proxyport=8080;no_verify_ssl=true;storage_class=STANDARD;part_size=7340032;restore_days=3")
	if err != nil {
		t.Fatalf("ParseConfig(): %v", err)
	}
	if cfg.Region != "cn-east-1" || cfg.ProxyHost != "proxy.local" || cfg.ProxyPort != "8080" {
		t.Fatalf("unexpected compatibility fields: %#v", cfg)
	}
	if !cfg.NoVerifySSL || cfg.StorageClass != "STANDARD" || cfg.PartSize != 7340032 || cfg.RestoreDays != 3 {
		t.Fatalf("unexpected sio compatibility settings: %#v", cfg)
	}
}

func TestBuildHTTPClientHonorsProxyAndTLSFlags(t *testing.T) {
	client := buildHTTPClient(Config{ProxyHost: "proxy.local", ProxyPort: "8080", NoVerifySSL: true})
	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("unexpected transport: %#v", client.Transport)
	}
	if transport.TLSClientConfig == nil || !transport.TLSClientConfig.InsecureSkipVerify {
		t.Fatalf("unexpected tls config: %#v", transport.TLSClientConfig)
	}
	proxyURL, err := transport.Proxy(&http.Request{})
	if err != nil {
		t.Fatalf("proxy(): %v", err)
	}
	if proxyURL.String() != "http://proxy.local:8080" {
		t.Fatalf("proxy url = %v", proxyURL)
	}
}

func TestNormalizeBucketNameForSIOProfiles(t *testing.T) {
	for _, backend := range []string{"sio", "siov1", "gdas"} {
		if got := normalizeBucketName(backend, "bucket/subdir"); got != "bucket" {
			t.Fatalf("%s normalized bucket = %q", backend, got)
		}
	}
	if got := normalizeBucketName("s3", "bucket/subdir"); got != "bucket/subdir" {
		t.Fatalf("s3 normalized bucket = %q", got)
	}
}

func TestDeleteMissingErrorsAreTolerated(t *testing.T) {
	for _, err := range []error{
		errors.New("NoSuchBucket"),
		errors.New("NoSuchKey"),
		errors.New("404 Not Found"),
	} {
		if !isDeleteMissingError(err) {
			t.Fatalf("expected delete-missing tolerance for %v", err)
		}
	}
}
