package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type EndpointType string

const (
	EndpointTypeMock EndpointType = "mock"
	EndpointTypeS3   EndpointType = "s3"
	EndpointTypeSIO  EndpointType = "sio"
)

type EndpointConfig struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Type        EndpointType `json:"type"`
	Endpoint    string       `json:"endpoint,omitempty"`
	Region      string       `json:"region,omitempty"`
	AccessKey   string       `json:"access_key,omitempty"`
	SecretKey   string       `json:"secret_key,omitempty"`
	PathStyle   bool         `json:"path_style"`
	ExtraConfig string       `json:"extra_config,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (c EndpointConfig) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("endpoint name cannot be empty")
	}
	switch c.Type {
	case EndpointTypeMock:
		return nil
	case EndpointTypeS3, EndpointTypeSIO:
		if strings.TrimSpace(c.Endpoint) == "" {
			return errors.New("endpoint address cannot be empty")
		}
		if strings.TrimSpace(c.AccessKey) == "" {
			return errors.New("access key cannot be empty")
		}
		if strings.TrimSpace(c.SecretKey) == "" {
			return errors.New("secret key cannot be empty")
		}
		return nil
	default:
		return fmt.Errorf("unsupported endpoint type: %q", c.Type)
	}
}

func (c EndpointConfig) RawConfig() string {
	parts := make([]string, 0, 6)
	appendKV := func(key, value string) {
		if strings.TrimSpace(value) == "" {
			return
		}
		parts = append(parts, fmt.Sprintf("%s=%s", key, strings.TrimSpace(value)))
	}
	appendKV("endpoint", c.Endpoint)
	appendKV("region", c.Region)
	appendKV("accesskey", c.AccessKey)
	appendKV("secretkey", c.SecretKey)
	if c.PathStyle {
		parts = append(parts, "path_style_access=true")
	}
	if extra := strings.TrimSpace(c.ExtraConfig); extra != "" {
		parts = append(parts, extra)
	}
	return strings.Join(parts, ";")
}
