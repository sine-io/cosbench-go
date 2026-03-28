package execution

import (
	"math/rand"
	"testing"
)

func TestParseOpConfig(t *testing.T) {
	pc, err := ParseOpConfig("cprefix=s3test;containers=r(1,2);objects=u(11,20);sizes=c(64)KiB;part_size=5242880;restore_days=3")
	if err != nil {
		t.Fatal(err)
	}
	if pc.PartSize != 5242880 {
		t.Fatalf("partSize = %d", pc.PartSize)
	}
	if pc.RestoreDays != 3 {
		t.Fatalf("restoreDays = %d", pc.RestoreDays)
	}
	r := rand.New(rand.NewSource(1))
	target := pc.NextTarget(r, 1, 1)
	if target.Bucket == "" || target.Key == "" {
		t.Fatalf("target = %#v", target)
	}
	if target.Size != 64*1024 {
		t.Fatalf("size = %d", target.Size)
	}
}

func TestParseOpConfigWithStorageFallback(t *testing.T) {
	pc, err := ParseOpConfigWithStorage("part_size=7340032;restore_days=3", "containers=c(1);objects=c(1)")
	if err != nil {
		t.Fatal(err)
	}
	if pc.PartSize != 7340032 {
		t.Fatalf("partSize = %d", pc.PartSize)
	}
	if pc.RestoreDays != 3 {
		t.Fatalf("restoreDays = %d", pc.RestoreDays)
	}
}

func TestParseOpConfigWithStorageOverride(t *testing.T) {
	pc, err := ParseOpConfigWithStorage("part_size=7340032;restore_days=3", "part_size=4194304;restore_days=5")
	if err != nil {
		t.Fatal(err)
	}
	if pc.PartSize != 4194304 {
		t.Fatalf("partSize = %d", pc.PartSize)
	}
	if pc.RestoreDays != 5 {
		t.Fatalf("restoreDays = %d", pc.RestoreDays)
	}
}

func TestParseOpConfigReadCompatibilityFlags(t *testing.T) {
	pc, err := ParseOpConfig("is_prefetch=true;is_range_request=true;file_length=15000000;chunk_length=5000000")
	if err != nil {
		t.Fatal(err)
	}
	if !pc.IsPrefetch {
		t.Fatal("expected prefetch flag")
	}
	if !pc.IsRangeRequest {
		t.Fatal("expected range-request flag")
	}
	if pc.FileLength != 15000000 {
		t.Fatalf("fileLength = %d", pc.FileLength)
	}
	if pc.ChunkLength != 5000000 {
		t.Fatalf("chunkLength = %d", pc.ChunkLength)
	}
}
