package trie

import (
	"errors"
	"strings"
)

var InvalidString = errors.New("InvalidString")

// TrieNode  trie Node
type TrieNode[T any] struct {
	children map[string]any
	Data     *T
	parent   any // pointer of parent node
	Key      string
}

// NewTrie create a new trie tree
func NewTrie[T any]() *TrieNode[T] {
	o := newNode[T]()
	return o
}

// newNode create a new trie node
func newNode[T any]() *TrieNode[T] {
	return &TrieNode[T]{
		children: map[string]any{},
	}
}

// NewChild create a child node of t
func (t *TrieNode[T]) NewChild() *TrieNode[T] {
	return &TrieNode[T]{
		children: map[string]any{},
		parent:   t,
	}
}

// Find walk through the tire and return the node that represent the key
// return nil if not found
func (t *TrieNode[T]) Find(key string) *TrieNode[T] {
	keySlice := strings.Split(key, "/")
	var pNode = t
	for _, lv := range keySlice {
		if _, ok := pNode.children[lv]; ok {
			pNode = pNode.children[lv].(*TrieNode[T])
		} else {
			return nil
		}
	}
	if pNode.Data != nil {
		return pNode
	}
	return nil
}

type OnMatch[T any] func(t *T) bool

// MatchKey walk through the tire and call the fn callback for each message witch match the key filter.
func (t *TrieNode[T]) MatchKey(keySlice []string, fn OnMatch[T]) {
	endFlag := len(keySlice) == 1
	switch keySlice[0] {
	case "#":
		t.preOrderTraverse(fn)
	case "+":
		// 当前层的所有
		for _, o := range t.children {
			v := o.(*TrieNode[T])
			if endFlag {
				if v.Data != nil {
					fn(v.Data)
				}
			} else {
				v.MatchKey(keySlice[1:], fn)
			}
		}
	default:
		if o := t.children[keySlice[0]]; o != nil {

			n := o.(*TrieNode[T])
			if endFlag {
				if n.Data != nil {
					fn(n.Data)
				}
			} else {
				n.MatchKey(keySlice[1:], fn)
			}
		}
	}
}

func (t *TrieNode[T]) GetMatchedData(filter string) []*T {
	keyLv := strings.Split(filter, "/")
	var rs []*T
	t.MatchKey(keyLv, func(d *T) bool {
		var o *T
		if d != nil {
			t1 := *d
			o = new(T)
			*o = t1
		}
		rs = append(rs, o)
		return true
	})
	return rs
}

// AddData add data
func (t *TrieNode[T]) AddData(key string, data *T) error {

	keySlice := strings.Split(key, "/")
	var pNode = t
	for _, lv := range keySlice {
		if _, ok := pNode.children[lv]; !ok {
			pNode.children[lv] = pNode.NewChild()
		}
		pNode = pNode.children[lv].(*TrieNode[T])
	}
	pNode.Data = data
	pNode.Key = key

	return nil
}

func (t *TrieNode[T]) Remove(key string) {
	keySlice := strings.Split(key, "/")
	l := len(keySlice)
	var pNode = t
	for _, lv := range keySlice {
		if _, ok := pNode.children[lv]; ok {
			pNode = pNode.children[lv].(*TrieNode[T])
		} else {
			return
		}
	}
	pNode.Data = nil
	if len(pNode.children) == 0 {
		delete(pNode.parent.(*TrieNode[T]).children, keySlice[l-1])
	}
}

func (t *TrieNode[T]) preOrderTraverse(fn OnMatch[T]) bool {
	if t == nil {
		return false
	}
	if t.Data != nil {
		if !fn(t.Data) {
			return false
		}
	}
	for _, c := range t.children {
		if !c.(*TrieNode[T]).preOrderTraverse(fn) {
			return false
		}
	}
	return true
}
