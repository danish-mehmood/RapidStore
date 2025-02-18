package btree

import (
	"strings"
	"sync"
)

const (
	degree = 4 // Minimum degree of B-tree
)

// Item represents a key-value pair
type Item struct {
	Key   string
	Value string
}

// Node represents a node in the B-tree
type Node struct {
	Items    []Item
	Children []*Node
	IsLeaf   bool
}

// BTree represents a B-tree data structure
type BTree struct {
	root *Node
	mu   sync.RWMutex
}

// NewBTree creates a new B-tree
func NewBTree() *BTree {
	return &BTree{
		root: &Node{
			Items:    make([]Item, 0),
			Children: make([]*Node, 0),
			IsLeaf:   true,
		},
	}
}

// search returns the index where the key should be inserted
func (n *Node) search(key string) int {
	for i := 0; i < len(n.Items); i++ {
		if strings.Compare(key, n.Items[i].Key) <= 0 {
			return i
		}
	}
	return len(n.Items)
}

// Get retrieves a value by key
func (t *BTree) Get(key string) (string, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	node := t.root
	for {
		i := node.search(key)
		if i < len(node.Items) && node.Items[i].Key == key {
			return node.Items[i].Value, true
		}
		if node.IsLeaf {
			return "", false
		}
		node = node.Children[i]
	}
}

// splitChild splits a full child node
func (n *Node) splitChild(i int) {
	child := n.Children[i]
	newChild := &Node{
		Items:    make([]Item, 0, 2*degree-1),
		Children: make([]*Node, 0, 2*degree),
		IsLeaf:   child.IsLeaf,
	}

	// Move items to new child
	medianIdx := degree - 1
	medianItem := child.Items[medianIdx]

	// Copy items after median to new child
	newChild.Items = append(newChild.Items, child.Items[medianIdx+1:]...)
	child.Items = child.Items[:medianIdx]

	// Handle children if not a leaf
	if !child.IsLeaf {
		newChild.Children = append(newChild.Children, child.Children[medianIdx+1:]...)
		child.Children = child.Children[:medianIdx+1]
	}

	// Insert new child into parent
	n.Items = append(n.Items, Item{})
	copy(n.Items[i+1:], n.Items[i:])
	n.Items[i] = medianItem

	n.Children = append(n.Children, nil)
	copy(n.Children[i+2:], n.Children[i+1:])
	n.Children[i+1] = newChild
}

// findKey searches for a key in a node and returns (index, found)
func (n *Node) findKey(key string) (int, bool) {
	for i, item := range n.Items {
		if item.Key == key {
			return i, true
		}
		if strings.Compare(key, item.Key) < 0 {
			return i, false
		}
	}
	return len(n.Items), false
}

// insertNonFull inserts a key-value pair into a non-full node
func (n *Node) insertNonFull(key string, value string) {
	i := len(n.Items) - 1

	if n.IsLeaf {
		// Check if key exists in leaf node
		if idx, found := n.findKey(key); found {
			// Update existing value
			n.Items[idx].Value = value
			return
		}
		// Insert into leaf node
		n.Items = append(n.Items, Item{})
		for i >= 0 && strings.Compare(key, n.Items[i].Key) < 0 {
			n.Items[i+1] = n.Items[i]
			i--
		}
		n.Items[i+1] = Item{Key: key, Value: value}
	} else {
		// Find child to recurse into
		i = n.search(key)

		// Check if key exists in current node
		if i < len(n.Items) && n.Items[i].Key == key {
			n.Items[i].Value = value
			return
		}

		if len(n.Children[i].Items) == 2*degree-1 {
			n.splitChild(i)
			if strings.Compare(key, n.Items[i].Key) > 0 {
				i++
			}
		}
		n.Children[i].insertNonFull(key, value)
	}
}

// Set inserts or updates a key-value pair
func (t *BTree) Set(key string, value string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	root := t.root
	if len(root.Items) == 2*degree-1 {
		newRoot := &Node{
			Items:    make([]Item, 0, 2*degree-1),
			Children: []*Node{root},
			IsLeaf:   false,
		}
		t.root = newRoot
		newRoot.splitChild(0)
		newRoot.insertNonFull(key, value)
	} else {
		root.insertNonFull(key, value)
	}
	return nil
}

// List returns all keys in sorted order
func (t *BTree) List() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var keys []string
	var traverse func(*Node)
	traverse = func(n *Node) {
		for i := 0; i < len(n.Items); i++ {
			if !n.IsLeaf {
				traverse(n.Children[i])
			}
			keys = append(keys, n.Items[i].Key)
		}
		if !n.IsLeaf {
			traverse(n.Children[len(n.Items)])
		}
	}
	traverse(t.root)
	return keys
}

// Delete removes a key-value pair
func (t *BTree) Delete(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	node := t.root
	var path []*Node
	var indices []int

	// Find the key and track the path
	for {
		i := node.search(key)
		if i < len(node.Items) && node.Items[i].Key == key {
			// Key found, store path and index
			path = append(path, node)
			indices = append(indices, i)
			break
		}
		if node.IsLeaf {
			return false // Key not found
		}
		path = append(path, node)
		indices = append(indices, i)
		node = node.Children[i]
	}

	if len(path) == 0 {
		return false
	}

	// Get the leaf node containing the key
	node = path[len(path)-1]
	keyIndex := indices[len(indices)-1]

	if node.IsLeaf {
		// Simple case: remove from leaf
		copy(node.Items[keyIndex:], node.Items[keyIndex+1:])
		node.Items = node.Items[:len(node.Items)-1]
		return true
	}

	// If not a leaf, replace with predecessor or successor
	predecessor := node.Children[keyIndex]
	for !predecessor.IsLeaf {
		predecessor = predecessor.Children[len(predecessor.Children)-1]
	}

	// Replace the key-value pair with the rightmost item from predecessor
	node.Items[keyIndex] = predecessor.Items[len(predecessor.Items)-1]
	predecessor.Items = predecessor.Items[:len(predecessor.Items)-1]

	return true
}
