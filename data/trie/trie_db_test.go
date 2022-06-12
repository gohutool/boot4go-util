package trie

import (
	"fmt"
	"testing"
)

type Sample struct {
	ID string
}

func TestTrie(t *testing.T) {

	trie := NewTrie[Sample]()
	trie.AddData("a/b/c1", &Sample{ID: "a/b/c1"})
	trie.AddData("a/c/c2", &Sample{ID: "a/c/c2"})
	trie.AddData("a/d/c3", &Sample{ID: "a/d/c3"})

	trie.AddData("/a/b/c1", &Sample{ID: "/a/b/c1"})
	trie.AddData("/a/c/c2", &Sample{ID: "/a/b/c2"})
	trie.AddData("/a/d/c3", &Sample{ID: "/a/b/c3"})
	trie.AddData("/a/b/c2", &Sample{ID: "/a/b/c2"})
	trie.AddData("/a/ccc", &Sample{ID: "/a/ccc"})
	trie.AddData("/a/ddd", &Sample{ID: "/a/ddd"})

	r := trie.GetMatchedData("a/#")
	printSample(r)

	r = trie.GetMatchedData("/a/#")
	printSample(r)

	r = trie.GetMatchedData("/a/b/#")
	printSample(r)

	r = trie.GetMatchedData("/a/+")
	printSample(r)
}

func printSample(l []*Sample) {

	for _, sample := range l {
		fmt.Printf("%+v ", *sample)
	}
	fmt.Println()
}
