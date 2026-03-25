package execution

import (
	"math/rand"
	"testing"

	"github.com/sine-io/cosbench-go/internal/domain/workload"
)

func TestWeightedOperationPicker(t *testing.T) {
	p, err := NewWeightedOperationPicker([]workload.Operation{{Type: "read", Ratio: 80}, {Type: "write", Ratio: 20}})
	if err != nil {
		t.Fatal(err)
	}
	r := rand.New(rand.NewSource(1))
	got := p.Pick(r)
	if got.Type == "" {
		t.Fatal("empty op")
	}
}
