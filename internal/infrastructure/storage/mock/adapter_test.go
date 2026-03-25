package mock

import (
	"context"
	"strings"
	"testing"
)

func TestMockAdapterBasicFlow(t *testing.T) {
	a := New()
	ctx := context.Background()
	if err := a.CreateBucket(ctx, "b1"); err != nil {
		t.Fatal(err)
	}
	if err := a.PutObject(ctx, "b1", "k1", strings.NewReader("abc"), 3); err != nil {
		t.Fatal(err)
	}
	meta, err := a.HeadObject(ctx, "b1", "k1")
	if err != nil {
		t.Fatal(err)
	}
	if meta.ContentLength != 3 {
		t.Fatalf("content length = %d", meta.ContentLength)
	}
	items, err := a.ListObjects(ctx, "b1", "k", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("items = %d", len(items))
	}
}
