package execution

import (
	"math/rand"
	"testing"
)

func TestParseIntGenerator(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	g, err := ParseIntGenerator("c(64)")
	if err != nil {
		t.Fatal(err)
	}
	if got := g.Next(r, 1, 1); got != 64 {
		t.Fatalf("constant = %d", got)
	}

	g, err = ParseIntGenerator("s(1,3)")
	if err != nil {
		t.Fatal(err)
	}
	if g.Next(r, 1, 1) != 1 || g.Next(r, 1, 1) != 2 || g.Next(r, 1, 1) != 3 {
		t.Fatal("sequential generator broken")
	}

	g, err = ParseIntGenerator("r(1,3)")
	if err != nil {
		t.Fatal(err)
	}
	if got := g.Next(r, 1, 1); got < 1 || got > 3 {
		t.Fatalf("range = %d", got)
	}

	g, err = ParseIntGenerator("r(0,0)")
	if err != nil {
		t.Fatal(err)
	}
	if got := g.Next(r, 1, 1); got != 0 {
		t.Fatalf("zero range = %d", got)
	}
}
