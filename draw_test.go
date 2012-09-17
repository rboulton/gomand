package main

import "testing"

func TestIterate(t *testing.T) {
	x, y, maxiters, expected := 0.0, 0.0, 10, 0
	got := iterate(x, y, maxiters)
	if got != expected {
		t.Errorf("iterate(%f, %f, %d) = %d, want %d", x, y, maxiters, got, expected)
	}
}
