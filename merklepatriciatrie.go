package merklepatriciatrie

import (
	"errors"
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

type Node struct {
	isLeaf   bool
	value    []byte
	children map[byte]*Node
}

type MPT struct {
	root *Node
}

func NewMPT() *MPT {
	return &MPT{
		root: &Node{children: make(map[byte]*Node)},
	}
}

func (t *MPT) Get(key []byte) ([]byte, error) {
	currentNode := t.root

	// Node look up by key
	for _, byteChar := range key {
		// Check if the current key exists, if not return an error
		if _, byteFound := currentNode.children[byteChar]; !byteFound {
			return nil, errors.New("key not found")
		}

		currentNode = currentNode.children[byteChar]

	}

	// Return the node actual value only if it's a leaf
	if currentNode.isLeaf {
		return currentNode.value, nil
	}

	// Otherwise return an error
	return nil, errors.New("key found but node is not a leaf")
}

func (t *MPT) Put(key []byte, value []byte) {
	currentNode := t.root

	for _, byteChar := range key {
		// Check if the key already exists
		if _, byteFound := currentNode.children[byteChar]; !byteFound {
			// When key is not found, we create a new child
			currentNode.children[byteChar] = &Node{children: make(map[byte]*Node)}
		}

		// when we find the corresponding key, then we swich node to the next one
		currentNode = currentNode.children[byteChar]
	}

	// Assign value to current node
	currentNode.value = value
	// Assign new leaf indicator
	currentNode.isLeaf = true

}
