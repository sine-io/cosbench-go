package app

import "testing"

func TestNewApp(t *testing.T) {
	application, err := New(Config{DataDir: t.TempDir(), ViewDir: "../../web/templates"})
	if err != nil {
		t.Fatalf("New(): %v", err)
	}
	if application.Handler == nil || application.Manager == nil {
		t.Fatalf("unexpected app: %#v", application)
	}
}
