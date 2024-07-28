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

func TestNonExistingKey(t *testing.T) {
	mpt := NewMPT()

	mpt.Put([]byte("key1"), []byte("value1"))
	_, got := mpt.Get([]byte("key2"))
	want := "key not found"

	if got == nil {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestDeletedElementShouldNotBeRetrivable(t *testing.T) {
	mpt := NewMPT()
	// Ensure the element exists
	mpt.Put([]byte("key1"), []byte("value1"))
	got, err := mpt.Get([]byte("key1"))
	want := []byte("value1")
	if err != nil {
		t.Errorf("go error %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %q want %q", got, want)
	}

	// Delete the element
	err = mpt.Del([]byte("key1"))

	if err != nil {
		t.Errorf("error when deleting element %q", err)
	}

	_, err = mpt.Get([]byte("key1"))
	if err == nil {
		t.Errorf("we should got an error when getting a deleted elements")
	}
}

func TestDeleteMultipleLevelsNode(t *testing.T) {
	mpt := NewMPT()
	// Ensure the element exists
	mpt.Put([]byte("key1"), []byte("value1"))
	mpt.Put([]byte("key2"), []byte("value3"))
	mpt.Put([]byte("key3"), []byte("value3"))

	// Delete the element
	err := mpt.Del([]byte("key3"))

	if err != nil {
		t.Errorf("error when deleting element %q", err)
	}

	// Getting key 1 and 2 should ok but should not key 3
	_, errKey1 := mpt.Get([]byte("key1"))
	_, errKey2 := mpt.Get([]byte("key2"))
	_, errKey3 := mpt.Get([]byte("key3"))

	if errKey1 != nil || errKey2 != nil {
		t.Errorf("we should not get an error when getting these elements")
	}

	if errKey3 == nil {
		t.Errorf("we should got an error when getting a deleted elements")
	}

}

func TestCommit(t *testing.T) {
	mpt := NewMPT()
	mpt.Put([]byte("key1"), []byte("value1"))
	mpt.Put([]byte("key2"), []byte("value2"))
	rootHash := mpt.Commit()
	if len(rootHash) == 0 {
		t.Fatalf("expected a valid root hash, got empty hash")
	}
}

func TestProof(t *testing.T) {
	mpt := NewMPT()
	mpt.Put([]byte("key1"), []byte("value1"))
	mpt.Commit()
	proof, err := mpt.Proof([]byte("key1"))
	if err != nil {
		t.Fatalf("expected proof for key1, got error: %v", err)
	}

	if len(proof) == 0 {
		t.Fatalf("expected a valid proof, got empty proof")
	}

}
func TestProofWithNonExistingKey(t *testing.T) {
	mpt := NewMPT()
	mpt.Put([]byte("key1"), []byte("value1"))
	mpt.Commit()
	_, err := mpt.Proof([]byte("key2"))
	if err == nil {
		t.Fatalf("expected error for key2 proof, got none")
	}

}
