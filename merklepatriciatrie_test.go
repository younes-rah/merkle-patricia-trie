package merklepatriciatrie

import (
	"bytes"
	"testing"
)

func TestEmptyMptGet(t *testing.T) {
	mpt := NewMPT()
	_, got := mpt.Get([]byte("key1"))
	want := "key not found"

	if got == nil {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestInsertedElementIsRetrivable(t *testing.T) {
	mpt := NewMPT()

	mpt.Put([]byte("key1"), []byte("value1"))
	got, err := mpt.Get([]byte("key1"))
	want := []byte("value1")
	if err != nil {
		t.Errorf("go error %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %q want %q", got, want)
	}
}
