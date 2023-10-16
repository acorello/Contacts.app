package seq

import (
	"math"
	"slices"
	"testing"
)

func TestMap(t *testing.T) {
	res := Map(math.Abs)
	if res != nil {
		t.Errorf("Expected nil but got %v", res)
	}
	res = Map(math.Abs, 1, -1)
	expected := []float64{1, 1}
	if !slices.Equal(res, expected) {
		t.Errorf("expected %#v but got %#v", expected, res)
	}
}
