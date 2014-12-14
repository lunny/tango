package tango

import "testing"

func TestIsNil(t *testing.T) {
	if !isNil(nil) {
		t.Error("nil")
	}

	if isNil(1) {
		t.Error("1")
	}

	if isNil("tttt") {
		t.Error("tttt")
	}

	type A struct {
	}

	var a A

	if isNil(a) {
		t.Error("a0")
	}

	if isNil(&a) {
		t.Error("a")
	}

	if isNil(new(A)) {
		t.Error("a2")
	}

	var b *A
	if !isNil(b) {
		t.Error("b")
	}

	var c interface{}
	if !isNil(c) {
		t.Error("c")
	}
}
