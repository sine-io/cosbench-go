package config

import "testing"

func TestParseKVConfig(t *testing.T) {
	got := ParseKVConfig(" accesskey = ak ; secretkey= sk ; endpoint = http://127.0.0.1 ; invalid ")
	if got["accesskey"] != "ak" {
		t.Fatalf("accesskey = %q", got["accesskey"])
	}
	if got["secretkey"] != "sk" {
		t.Fatalf("secretkey = %q", got["secretkey"])
	}
	if got["endpoint"] != "http://127.0.0.1" {
		t.Fatalf("endpoint = %q", got["endpoint"])
	}
	if _, ok := got["invalid"]; ok {
		t.Fatal("invalid entry should be ignored")
	}
}
