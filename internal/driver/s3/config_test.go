package s3

import (
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
