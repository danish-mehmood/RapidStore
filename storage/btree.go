package storage

import (
	"github.com/danish-mehmood/RapidStore/btree"
	"github.com/danish-mehmood/RapidStore/errors"
)

// BTreeEngine implements the Engine interface using a B-tree
type BTreeEngine struct {
	tree *btree.BTree
}

// NewBTreeEngine creates a new BTreeEngine instance
func NewBTreeEngine() *BTreeEngine {

	return &BTreeEngine{
		tree: btree.NewBTree(),
	}
}

// Get retrieves a value by key
func (e *BTreeEngine) Get(key string) (string, error) {
	if key == "" {
		return "", errors.ErrEmptyKey
	}

	if value, exists := e.tree.Get(key); exists {
		return value, nil
	}
	return "", errors.ErrKeyNotFound
}

// Set stores a key-value pair
func (e *BTreeEngine) Set(key, value string) error {
	if key == "" {
		return errors.ErrEmptyKey
	}
	return e.tree.Set(key, value)
}

// Delete removes a key-value pair
func (e *BTreeEngine) Delete(key string) error {
	if key == "" {
		return errors.ErrEmptyKey
	}

	if !e.tree.Delete(key) {
		return errors.ErrKeyNotFound
	}
	return nil
}

// List returns all keys in sorted order
func (e *BTreeEngine) List() []string {
	return e.tree.List()
}
