package merklepatriciatriev2

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutAndGet(t *testing.T) {
	trie := NewMPT()

	// Test Put and Get
	trie.Put([]byte("key1"), []byte("value1"))
	value, err := trie.Get([]byte("key1"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("value1"), value)

	// Test Put and Get another key
	trie.Put([]byte("key2"), []byte("value2"))
	value, err = trie.Get([]byte("key2"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("value2"), value)
}

func TestDelete(t *testing.T) {
	trie := NewMPT()

	// Put a key-value pair
	trie.Put([]byte("key1"), []byte("value1"))

	// Test Del
	err := trie.Del([]byte("key1"))
	assert.NoError(t, err)
	_, err = trie.Get([]byte("key1"))
	assert.Error(t, err)
}

func TestCommit(t *testing.T) {
	trie := NewMPT()

	// Put a key-value pair
	trie.Put([]byte("key1"), []byte("value1"))

	// Test Commit
	rootKey := trie.Commit()
	assert.NotNil(t, rootKey)

	// Verify the trie is saved in persistent storage
	storedRootKey, exists := storage[hex.EncodeToString(rootKey)]
	assert.True(t, exists)
	assert.Equal(t, rootKey, storedRootKey)
}

func TestProofSize(t *testing.T) {
	trie := NewMPT()

	// Populate the trie with two key-value pairs
	keys := []string{"key1", "key2"}
	values := []string{"value1", "value2"}

	for i := range keys {
		trie.Put(keyToNibbles([]byte(keys[i])), []byte(values[i]))
	}

	// Generate proofs for the keys
	proof, err := trie.Proof(keyToNibbles([]byte("key1")))
	assert.NoError(t, err)
	assert.NotEmpty(t, proof)

	// Verify the size of the proof
	// Should be 3 because in this setup we expect:
	// 1 ExtensionNode since key1 and key2 share the prefix "key"
	// 1 BranchNode the follow the first extention node
	// 2 LeafNode, for key1 and key2

	assert.Equal(t, 3, len(proof), "unexpected proof size: %d", len(proof))
}
func TestEthereumLikeData(t *testing.T) {
	trie := NewMPT()

	ethereumData := map[string]string{
		"0x0000000000000000000000000000000000000001": "0x1000000000000000000000000000000000000000",
		"0x0000000000000000000000000000000000000002": "0x2000000000000000000000000000000000000000",
		"0x0000000000000000000000000000000000000003": "0x3000000000000000000000000000000000000000",
		"0x0000000000000000000000000000000000000004": "0x4000000000000000000000000000000000000000",
		"0x0000000000000000000000000000000000000005": "0x5000000000000000000000000000000000000000",
		"0x0000000000000000000000000000000000000006": "0x6000000000000000000000000000000000000000",
		"0x0000000000000000000000000000000000000007": "0x7000000000000000000000000000000000000000",
		"0x0000000000000000000000000000000000000008": "0x8000000000000000000000000000000000000000",
		"0x0000000000000000000000000000000000000009": "0x9000000000000000000000000000000000000000",
		"0x000000000000000000000000000000000000000a": "0xa000000000000000000000000000000000000000",
		"0x000000000000000000000000000000000000000b": "0xb000000000000000000000000000000000000000",
		"0x000000000000000000000000000000000000000c": "0xc000000000000000000000000000000000000000",
		"0x000000000000000000000000000000000000000d": "0xd000000000000000000000000000000000000000",
		"0x000000000000000000000000000000000000000e": "0xe000000000000000000000000000000000000000",
		"0x000000000000000000000000000000000000000f": "0xf000000000000000000000000000000000000000",
		"0x0000000000000000000000000000000000000010": "0x10000000000000000000000000000000000000000",
		"0x0000000000000000000000000000000000000011": "0x11000000000000000000000000000000000000000",
		"0x0000000000000000000000000000000000000012": "0x12000000000000000000000000000000000000000",
		"0x0000000000000000000000000000000000000013": "0x13000000000000000000000000000000000000000",
		"0x0000000000000000000000000000000000000014": "0x14000000000000000000000000000000000000000",
	}

	// Add all key-value pairs to the trie
	for key, value := range ethereumData {
		keyBytes, _ := hex.DecodeString(key[2:])
		valueBytes, _ := hex.DecodeString(value[2:])
		trie.Put(keyBytes, valueBytes)
	}

	// Verify all key-value pairs can be retrieved correctly
	for key, value := range ethereumData {
		keyBytes, _ := hex.DecodeString(key[2:])
		valueBytes, _ := hex.DecodeString(value[2:])
		retrievedValue, err := trie.Get(keyBytes)
		assert.NoError(t, err)
		assert.Equal(t, valueBytes, retrievedValue)
	}
}
func TestIntegration(t *testing.T) {
	trie := NewMPT()

	// Put a key-value pair
	trie.Put([]byte("key1"), []byte("value1"))

	// Get the value
	value, err := trie.Get([]byte("key1"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("value1"), value)

	// Delete the key-value pair
	err = trie.Del([]byte("key1"))
	assert.NoError(t, err)

	// Verify the key is deleted
	_, err = trie.Get([]byte("key1"))
	assert.Error(t, err)

	// Put another key-value pair to test Commit
	trie.Put([]byte("key2"), []byte("value2"))

	// Commit and verify the root key
	rootKey := trie.Commit()
	assert.NotNil(t, rootKey)

	// Verify the trie is saved in persistent storage
	storedRootKey, exists := storage[hex.EncodeToString(rootKey)]
	assert.True(t, exists)
	assert.Equal(t, rootKey, storedRootKey)

	// Test Proof
	proof, err := trie.Proof([]byte("key2"))
	assert.NoError(t, err)
	assert.NotEmpty(t, proof)
}
