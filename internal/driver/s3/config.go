package s3

import (
	"fmt"
	"strings"

	appconfig "github.com/sine-io/cosbench-go/internal/infrastructure/config"
)

type Config struct {
	Backend   string
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
	PathStyle bool
	ProxyHost string
	ProxyPort string
	NoVerifySSL bool
	StorageClass string
	PartSize int64
	RestoreDays int
	Raw       map[string]string
}

func ParseConfig(backend, raw string) (Config, error) {
	return ParseConfigMap(backend, appconfig.ParseKVConfig(raw))
}

func ParseConfigMap(backend string, values map[string]string) (Config, error) {
	kind := strings.ToLower(strings.TrimSpace(backend))
	if kind == "" {
		kind = "s3"
	}
	cfg := Config{
		Backend:   kind,
		Endpoint:  strings.TrimSpace(values["endpoint"]),
		Region:    firstNonEmpty(values["region"], values["aws_region"]),
		AccessKey: strings.TrimSpace(values["accesskey"]),
		SecretKey: strings.TrimSpace(values["secretkey"]),
		PathStyle: parseBool(values["path_style_access"]),
		ProxyHost: strings.TrimSpace(values["proxyhost"]),
		ProxyPort: strings.TrimSpace(values["proxyport"]),
		NoVerifySSL: parseBool(values["no_verify_ssl"]),
		StorageClass: strings.TrimSpace(values["storage_class"]),
		PartSize: parseInt64(values["part_size"], 5*1024*1024),
		RestoreDays: parseInt(values["restore_days"], 1),
		Raw:       values,
	}
	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}
	if kind == "sio" && values["path_style_access"] == "" {
		cfg.PathStyle = true
	}
	if cfg.Endpoint == "" {
		return Config{}, fmt.Errorf("%s endpoint is required", kind)
	}
	if cfg.AccessKey == "" {
		return Config{}, fmt.Errorf("%s access key is required", kind)
	}
	if cfg.SecretKey == "" {
		return Config{}, fmt.Errorf("%s secret key is required", kind)
	}
	return cfg, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func parseInt64(v string, def int64) int64 {
	v = strings.TrimSpace(v)
	if v == "" {
		return def
	}
	var n int64
	_, err := fmt.Sscan(v, &n)
	if err != nil {
		return def
	}
	return n
}

func parseInt(v string, def int) int {
	v = strings.TrimSpace(v)
	if v == "" {
		return def
	}
	var n int
	_, err := fmt.Sscan(v, &n)
	if err != nil {
		return def
	}
	return n
}

func parseBool(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
