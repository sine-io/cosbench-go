package execution

import (
	"fmt"
	"math/rand"

	"github.com/sine-io/cosbench-go/internal/domain/workload"
)

type WeightedOperationPicker struct {
	last int
	ops  []weightedOp
}

type weightedOp struct {
	upper int
	op    workload.Operation
}

func NewWeightedOperationPicker(ops []workload.Operation) (*WeightedOperationPicker, error) {
	p := &WeightedOperationPicker{}
	for _, op := range ops {
		if op.Ratio <= 0 {
			continue
		}
		p.last += op.Ratio
		p.ops = append(p.ops, weightedOp{upper: p.last, op: op})
	}
	if p.last != 100 {
		return nil, fmt.Errorf("op ratio should sum to 100, got %d", p.last)
	}
	return p, nil
}

func (p *WeightedOperationPicker) Pick(r *rand.Rand) workload.Operation {
	x := r.Intn(100) + 1
	for _, item := range p.ops {
		if x <= item.upper {
			return item.op
		}
	}
	panic("weighted picker invariant broken")
}
