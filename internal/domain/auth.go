package domain

type AuthSpec struct {
	Type   string `json:"type"`
	Config string `json:"config,omitempty"`
}
