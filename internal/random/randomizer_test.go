package random

import (
	"fmt"
	"testing"
)

func TestNewRandomizer(t *testing.T) {
	r := NewRandomizer()
	if r == nil {
		t.Error("NewRandomizer() returned nil")
	}

	if r.rnd == nil {
		t.Error("Randomizer.rnd is nil")
	}

	randomStr := r.String(10)
	fmt.Println(randomStr)
	if len(randomStr) != 10 {
		t.Errorf("Random string length is %d, expected 10", len(randomStr))
	}
}
