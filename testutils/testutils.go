package testutils

import (
	"bytes"
	"reflect"
	"testing"
)

func AssertEquals[T comparable](t *testing.T, got T, want T) {
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func AssertFalse(t *testing.T, got bool) {
	if got {
		t.Errorf("got true, expected false")
	}
}

func AssertTrue(t *testing.T, got bool) {
	if !got {
		t.Errorf("got false, expected true")
	}
}

func AssertMapEquals[T comparable](t *testing.T, got map[T]T, want map[T]T) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func AssertByteEquals(t *testing.T, got []byte, want []byte) {
	if !bytes.Equal(got, want) {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func AssertEmpty[T any](t *testing.T, slice []T) {
	if len(slice) != 0 {
		t.Errorf("got non empty slice")
	}
}

func AssertLen[T any](t *testing.T, slice []T, length int) {
	if len(slice) != length {
		t.Errorf("got a slice of len %v, wanted %v", len(slice), length)
	}
}
