package merklepatriciatriev2

import (
	"bytes"
	"crypto/sha256"
	"errors"

	"github.com/dgraph-io/badger/v3"
)

type iTrie interface {
	// Get returns the value associate with the key // error is returned if the key is not found.
	Get(key []byte) ([]byte, error)
	// Put inserts the [key,value] node in the trie
	Put(key []byte, value []byte)
	// Del removes a node from the trie
	// returns an error if not found.
	Del(key []byte) error
	// Commit saves the trie in persistent storage // and returns the trie root key.
	Commit() []byte
	// Proof returns the Merkle-proof associated with
	// a node. An error is returned if the node is not found.
	Proof(key []byte) ([][]byte, error)
}

var _ iTrie = (*MPT)(nil) // Ensures MPT implements iTrie

type NodeType int

const (
	BranchNode NodeType = iota
	ExtensionNode
	LeafNode
)

type Node struct {
	Type     NodeType
	Children [16]*Node
	Value    []byte
	Key      []byte
	Next     *Node
}

type Storage interface {
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
	Close() error
}

type InMemoryStorage struct {
	data map[string][]byte
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{data: make(map[string][]byte)}
}

func (s *InMemoryStorage) Put(key []byte, value []byte) error {
	s.data[string(key)] = value
	return nil
}

func (s *InMemoryStorage) Get(key []byte) ([]byte, error) {
	if value, ok := s.data[string(key)]; ok {
		return value, nil
	}
	return nil, errors.New("key not found")
}

func (s *InMemoryStorage) Delete(key []byte) error {
	delete(s.data, string(key))
	return nil
}

func (s *InMemoryStorage) Close() error {
	return nil
}

type BadgerStorage struct {
	db *badger.DB
}

func NewBadgerStorage(dbPath string) (*BadgerStorage, error) {
	opts := badger.DefaultOptions(dbPath)
	// Disable logging for simplicity
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &BadgerStorage{db: db}, nil
}

func (s *BadgerStorage) Put(key []byte, value []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func (s *BadgerStorage) Get(key []byte) ([]byte, error) {
	var value []byte
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		value, err = item.ValueCopy(nil)
		return err
	})
	return value, err
}

func (s *BadgerStorage) Delete(key []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func (s *BadgerStorage) Close() error {
	return s.db.Close()
}

type MPT struct {
	root    *Node
	storage Storage
}

func NewMPT(storage Storage) *MPT {
	return &MPT{storage: storage}
}

func (t *MPT) Get(key []byte) ([]byte, error) {
	node := t.root
	nibbles := keyToNibbles(key)

	for node != nil {
		switch node.Type {
		case BranchNode:
			if len(nibbles) == 0 {
				return node.Value, nil
			}
			node = node.Children[nibbles[0]]
			nibbles = nibbles[1:]
		case ExtensionNode:
			prefix := node.Key
			if !bytes.HasPrefix(nibbles, prefix) {
				return nil, errors.New("key not found")
			}
			nibbles = nibbles[len(prefix):]
			node = node.Next
		case LeafNode:
			if bytes.Equal(nibbles, node.Key) {
				return node.Value, nil
			}
			return nil, errors.New("key not found")
		}
	}

	return nil, errors.New("key not found")
}

func keyToNibbles(key []byte) []byte {
	nibbles := make([]byte, len(key)*2)
	for i, b := range key {
		nibbles[i*2] = b / 16
		nibbles[i*2+1] = b % 16
	}
	return nibbles
}

func (t *MPT) Put(key []byte, value []byte) {
	nibbles := keyToNibbles(key)
	t.root = t.put(t.root, nibbles, value)
}

func (t *MPT) put(node *Node, nibbles []byte, value []byte) *Node {
	if node == nil {
		return &Node{Type: LeafNode, Key: nibbles, Value: value}
	}

	switch node.Type {
	case BranchNode:
		if len(nibbles) == 0 {
			node.Value = value
			return node
		}
		node.Children[nibbles[0]] = t.put(node.Children[nibbles[0]], nibbles[1:], value)
	case ExtensionNode:
		prefix := node.Key
		commonPrefixLen := commonPrefixLength(nibbles, prefix)
		if commonPrefixLen == len(prefix) {
			nibbles = nibbles[commonPrefixLen:]
			node.Next = t.put(node.Next, nibbles, value)
			return node
		}
		newBranch := &Node{Type: BranchNode}
		newBranch.Children[prefix[commonPrefixLen]] = node.Next
		node.Next = newBranch
		node.Key = prefix[:commonPrefixLen]
		nibbles = nibbles[commonPrefixLen:]
		if len(nibbles) == 0 {
			newBranch.Value = value
		} else {
			newBranch.Children[nibbles[0]] = t.put(newBranch.Children[nibbles[0]], nibbles[1:], value)
		}
		return node
	case LeafNode:
		if bytes.Equal(nibbles, node.Key) {
			node.Value = value
			return node
		}
		commonPrefixLen := commonPrefixLength(nibbles, node.Key)
		newBranch := &Node{Type: BranchNode}
		if commonPrefixLen == len(node.Key) {
			newBranch.Value = node.Value
		} else {
			newBranch.Children[node.Key[commonPrefixLen]] = &Node{Type: LeafNode, Key: node.Key[commonPrefixLen+1:], Value: node.Value}
		}
		if commonPrefixLen == len(nibbles) {
			newBranch.Value = value
		} else {
			newBranch.Children[nibbles[commonPrefixLen]] = t.put(newBranch.Children[nibbles[commonPrefixLen]], nibbles[commonPrefixLen+1:], value)
		}
		if commonPrefixLen == 0 {
			return newBranch
		}
		return &Node{Type: ExtensionNode, Key: nibbles[:commonPrefixLen], Next: newBranch}
	}

	return node
}

func commonPrefixLength(a, b []byte) int {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}
	for i := 0; i < minLen; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return minLen
}

func (t *MPT) Del(key []byte) error {
	nibbles := keyToNibbles(key)
	var err error
	t.root, err = t.del(t.root, nibbles)
	return err
}

func (t *MPT) del(node *Node, nibbles []byte) (*Node, error) {
	if node == nil {
		return nil, errors.New("key not found")
	}

	switch node.Type {
	case BranchNode:
		if len(nibbles) == 0 {
			node.Value = nil
			return node, nil
		}
		var err error
		node.Children[nibbles[0]], err = t.del(node.Children[nibbles[0]], nibbles[1:])
		if err != nil {
			return nil, err
		}
	case ExtensionNode:
		prefix := node.Key
		if !bytes.HasPrefix(nibbles, prefix) {
			return nil, errors.New("key not found")
		}
		nibbles = nibbles[len(prefix):]
		var err error
		node.Next, err = t.del(node.Next, nibbles)
		if err != nil {
			return nil, err
		}
	case LeafNode:
		if bytes.Equal(nibbles, node.Key) {
			return nil, nil
		}
		return nil, errors.New("key not found")
	}

	return node, nil
}

func (t *MPT) Commit() []byte {
	if t.root == nil {
		return nil
	}
	rootHash := t.hashNode(t.root)
	t.storage.Put(rootHash, rootHash)

	return rootHash
}

func (t *MPT) hashNode(node *Node) []byte {
	if node == nil {
		return nil
	}

	switch node.Type {
	case BranchNode:
		hashes := make([]byte, 32*17)
		for i := 0; i < 16; i++ {
			childHash := t.hashNode(node.Children[i])
			copy(hashes[i*32:(i+1)*32], childHash)
		}
		if node.Value != nil {
			copy(hashes[16*32:], node.Value)
		}
		return hashToSlice(sha256.Sum256(hashes))
	case ExtensionNode:
		hashedNext := t.hashNode(node.Next)
		combined := append(node.Key, hashedNext...)
		return hashToSlice(sha256.Sum256(combined))
	case LeafNode:
		combined := append(node.Key, node.Value...)
		return hashToSlice(sha256.Sum256(combined))
	}

	return nil
}

func hashToSlice(hash [32]byte) []byte {
	return hash[:]
}

func (t *MPT) Proof(key []byte) ([][]byte, error) {
	nibbles := keyToNibbles(key)
	proof := [][]byte{}
	node := t.root

	for node != nil {
		switch node.Type {
		case BranchNode:
			if len(nibbles) == 0 {
				return proof, nil
			}
			proof = append(proof, t.hashNode(node))
			node = node.Children[nibbles[0]]
			nibbles = nibbles[1:]
		case ExtensionNode:
			prefix := node.Key
			if !bytes.HasPrefix(nibbles, prefix) {
				return nil, errors.New("key not found")
			}
			proof = append(proof, t.hashNode(node))
			nibbles = nibbles[len(prefix):]
			node = node.Next
		case LeafNode:
			if bytes.Equal(nibbles, node.Key) {
				proof = append(proof, t.hashNode(node))
				return proof, nil
			}
			return nil, errors.New("key not found")
		}
	}

	return nil, errors.New("key not found")
}
