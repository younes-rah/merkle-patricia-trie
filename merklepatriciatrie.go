package merklepatriciatrie

import (
	"crypto/sha256"
	"errors"

	"github.com/syndtr/goleveldb/leveldb"
)

const DB_FILE_NAME = "db"

type iTrie interface {
	// Get returns the value associate with the key
	// error is returned if the key is not found.
	Get(key []byte) ([]byte, error)
	// Put inserts the [key,value] node in the trie
	Put(key []byte, value []byte)
	// Del removes a node from the trie
	// returns an error if not found.
	Del(key []byte) error
	// Commit saves the trie in persistent storage
	// and returns the trie root key.
	Commit() []byte
	// Proof returns the Merkle-proof associated with
	// a node. An error is returned if the node is not found.
	Proof(key []byte) ([][]byte, error)
}

type Node struct {
	isLeaf   bool
	value    []byte
	children map[byte]*Node
	hash     []byte
}

type MPT struct {
	root    *Node
	storage *leveldb.DB
}

func NewMPT() *MPT {
	db, err := leveldb.OpenFile(DB_FILE_NAME, nil)

	if err != nil {
		panic("Unable to open database")
	}
	// defer db.Close()

	return &MPT{
		root:    &Node{children: make(map[byte]*Node)},
		storage: db,
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

func (t *MPT) Del(key []byte) error {
	// Buid a tree path "tracer"
	pathTracer := []*Node{}

	currentNode := t.root

	for _, byteChar := range key {
		// Check if the current key exists, if not return an error
		if _, byteFound := currentNode.children[byteChar]; !byteFound {
			return errors.New("key not found")
		}
		// save path
		pathTracer = append(pathTracer, currentNode)
		currentNode = currentNode.children[byteChar]
	}

	// Once we went through the tree, check that the last node is a leaf node
	// Otherwise, we have an issue
	if !currentNode.isLeaf {
		return errors.New("key found but node is not a leaf")
	}

	// Deleting current node information
	currentNode.value = nil
	currentNode.isLeaf = false

	// Deleting node and rebuild back the path using our path trace

	for i := len(key) - 1; i >= 0; i-- {
		// init current node and byte char
		node := pathTracer[i]
		char := key[i]

		// ensure node is not a leaf and does not have children
		// before deleting from map
		// If that's the case we have finnished the deletion
		if currentNode.isLeaf || len(currentNode.children) > 0 {
			break
		}
		// The real delete is done
		delete(node.children, char)
		// Reassign the current node
		currentNode = node
	}

	// finnish with no error
	return nil

}

func (t *MPT) commitNode(node *Node) []byte {
	if node == nil {
		return nil
	}

	// Create hash
	nodeHash := sha256.New()

	// ensure current node is a leaf
	if node.isLeaf {
		nodeHash.Write(node.value)
	}

	// Append all childre hashes
	for char, child := range node.children {
		// We treat children recursively
		childHash := t.commitNode(child)
		nodeHash.Write([]byte{char})
		nodeHash.Write(childHash)

	}

	// Produce final hash
	node.hash = nodeHash.Sum(nil)
	t.storage.Put(node.hash, node.value, nil)

	return node.hash

}

func (t *MPT) Commit() []byte {
	t.commitNode(t.root)
	t.storage.Close()
	return t.root.hash
}

func (t *MPT) Proof(key []byte) ([][]byte, error) {
	currentNode := t.root

	proof := [][]byte{}

	for _, byteChar := range key {
		// Check if the key exists, if not return an error
		if _, byteCharFound := currentNode.children[byteChar]; !byteCharFound {
			return nil, errors.New("key not found")
		}
		proof = append(proof, currentNode.hash)
		currentNode = currentNode.children[byteChar]
	}

	// Ensure last element is a leaf node
	if currentNode.isLeaf {
		return proof, nil
	}

	// if that is not the case then we certainly have an issue in the structure
	return nil, errors.New("key found but node is not a leaf")

}
