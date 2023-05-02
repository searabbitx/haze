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
