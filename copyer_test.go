package copycat

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestCopyer(t *testing.T) {
	copyer := NewCopyer()
	if copyer == nil {
		t.Error("NewCopyer cannot return nil")
		return
	}
	if copyer.maxLevel != DefaultMaxLevel {
		t.Errorf("maxLevel should be equal to %v (current value: %v)", DefaultMaxLevel, copyer.maxLevel)
	}
	if copyer.flags != 0 {
		t.Errorf("flags should be equal to 0 (current value: %v)", copyer.flags)
	}
}

func TestGetSetFlags(t *testing.T) {
	copyer := NewCopyer()
	if copyer == nil {
		t.Error("NewCopyer cannot return nil")
		return
	}
	flagsToCheck := []Flags{
		FPreserveHierarchy, FCopyChan, FCopyFunc, FCopyUintptr, FCopyUnsafePointer, FCopyInterface, FAll,
	}
	for _, expected := range flagsToCheck {
		copyer.SetFlags(expected)
		if copyer.GetFlags() != expected {
			t.Errorf("Get/SetFlags: expected %v, got %v", expected, copyer.GetFlags())
		}
	}
}

func TestGetSetMaxLevel(t *testing.T) {
	copyer := NewCopyer()
	if copyer == nil {
		t.Error("NewCopyer cannot return nil")
		return
	}
	randomLevel := uint(rand.Uint32())
	copyer.SetMaxLevel(randomLevel)
	if copyer.GetMaxlevel() != randomLevel {
		t.Errorf("Get/SetMaxLevel: expected %v, got %v", randomLevel, copyer.GetMaxlevel())
	}
}

func TestCopyerDeepCopy(t *testing.T) {
	copyer := NewCopyer()
	if copyer == nil {
		t.Error("NewCopyer cannot return nil")
		return
	}

	type advance struct {
		A simple
		B simple
	}

	src := advance{
		A: simple{
			A: A,
			B: B,
			C: C,
			D: D,
			E: E,
			F: F,
			G: G,
			H: H,
			I: I,
			J: J,
			K: K,
			L: L,
		},
		B: simple{
			A: A,
			B: B,
			C: C,
			D: D,
			E: E,
			F: F,
			G: G,
			H: H,
			I: I,
			J: J,
			K: K,
			L: L,
		},
	}
	{
		dst := &advance{}
		expect := src
		copyer.DeepCopy(dst, src)
		if !reflect.DeepEqual(*dst, expect) {
			t.Errorf("dst != expect\ndst:\n%+v\nexpected:\n%+v", *dst, expect)
		}
	}
	{
		dst := &advance{}
		expect := advance{}
		copyer.SetMaxLevel(1)
		copyer.DeepCopy(dst, src)
		if !reflect.DeepEqual(*dst, expect) {
			t.Errorf("dst != expect\ndst:\n%+v\nexpected:\n%+v", *dst, expect)
		}
	}
}
