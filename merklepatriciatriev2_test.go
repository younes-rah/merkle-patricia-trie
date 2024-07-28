package merklepatriciatriev2

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutAndGet(t *testing.T) {
	trie := NewTrie()

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
	trie := NewTrie()

	// Put a key-value pair
	trie.Put([]byte("key1"), []byte("value1"))

	// Test Del
	err := trie.Del([]byte("key1"))
	assert.NoError(t, err)
	_, err = trie.Get([]byte("key1"))
	assert.Error(t, err)
}

func TestCommit(t *testing.T) {
	trie := NewTrie()

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

func TestProof(t *testing.T) {
	trie := NewTrie()

	// Test Proof (not implemented)
	_, err := trie.Proof([]byte("key2"))
	assert.Error(t, err)
}

func TestEthereumLikeData(t *testing.T) {
	trie := NewTrie()

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
	trie := NewTrie()

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

	// Test Proof (not implemented)
	_, err = trie.Proof([]byte("key2"))
	assert.Error(t, err)
}
