// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package trie

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

// ErrTrieLoadEmpty is thrown when you try and load an empty slice of strings
// to the trie
var ErrTrieLoadEmpty = errors.New("cannot load empty slice of strings ")

// Trie is a tree like data structure that allows us to process string finding
// operations faster than other means.
type Trie struct {
	root  *node
	count int
}

// New returns a new initialized trie
func New() *Trie {
	root := newNode(nil, rune(0))
	return &Trie{root, 0}
}

// Add adds a string to the trie creating any new nodes it needs
func (t *Trie) Add(s string) error {
	// fmt.Printf("string addded: %s\n", s)
	lower := strings.ToLower(s)
	rs := []rune(lower)

	if err := t.root.addChild(rs); err != nil {
		return err
	}
	t.count++
	return nil
}

// Load performs Add on a slice of strings.
func (t *Trie) Load(list []string) error {
	if len(list) == 0 {
		return ErrTrieLoadEmpty
	}

	for _, v := range list {
		if err := t.Add(v); err != nil {
			return err
		}
	}

	return nil
}

// LoadFile loads the contents of a json array of strings into the trie
func (t *Trie) LoadFile(name string) error {

	data, err := fileToStringSlice(name)

	if err != nil {
		return fmt.Errorf("erro converting file to []string: %s", err)
	}

	if err := t.Load(data); err != nil {
		return fmt.Errorf("error adding strings: %s", err)
	}

	return nil
}

func fileToStringSlice(name string) ([]string, error) {
	data := []string{}

	file, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("cannot read forbidden words file: %s", err)
	}

	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshall json into []string: %s", err)
	}

	return data, nil
}

// Find determines if an input string is exactly matches one present in
// the trie.
func (t *Trie) Find(s string) bool {
	ls := strings.ToLower(s)
	rs := []rune(ls)
	return t.root.isChild(rs)
}

// IsContained determins if there is a string in the trie contained within the
// input string. It also allows for a minimum length match.
func (t *Trie) IsContained(s string, min int) (bool, string) {
	ls := strings.ToLower(s)
	rs := []rune(ls)

	for i := range rs {
		result, sofar := t.root.isChildWithDepth(rs[i:], min, []rune(""))
		if result {
			return true, strings.TrimRight(string(sofar), "\x00")
		}
	}

	return false, ""
}

// Delete removes a string from the trie
func (t *Trie) Delete(s string) error {
	ls := strings.ToLower(s)
	rs := []rune(ls)
	if err := t.root.remove(rs); err != nil {
		return err
	}
	t.count--
	return nil
}

// Count returns the number of words in the trie
func (t *Trie) Count() int {
	return t.count
}

// Node is one item in a trie for computing relationships
type node struct {
	parent       *node
	children     map[rune]*node
	value        rune
	isTerminated bool
}

func newNode(parent *node, value rune) *node {
	children := make(map[rune]*node)
	return &node{parent, children, value, false}
}

func (n *node) addChild(value []rune) error {
	first, rest, _ := breakRuneSlice(value)
	ch, ok := n.children[first]
	if !ok {

		if len(value) == 0 {
			n.isTerminated = true
			return nil
		}

		ch = newNode(n, first)
		n.children[first] = ch

	}

	return ch.addChild(rest)
}

func (n *node) remove(value []rune) error {
	first, rest, _ := breakRuneSlice(value)

	if len(value) == 0 {
		n.isTerminated = false
		return nil
	}
	ch, ok := n.children[first]
	if ok {
		return ch.remove(rest)
	}

	return fmt.Errorf("could not find the children of node")

}

func breakRuneSlice(value []rune) (rune, []rune, rune) {
	first := rune(0)
	rest := []rune{}
	last := rune(0)

	if len(value) != 0 {
		first = value[0]
	}

	if len(value) > 1 {
		rest = value[1:]
		last = value[len(value)-1]
	}

	return first, rest, last
}

func (n *node) isChild(value []rune) bool {

	first, rest, _ := breakRuneSlice(value)

	ch, ok := n.children[first]
	if !ok {
		return false
	}
	if len(rest) == 0 {
		if ch.isTerminated {
			return true
		}
		return false
	}
	return ch.isChild(rest)

}

func (n *node) isChildWithDepth(value []rune, depth int, sofar []rune) (bool, []rune) {
	first, rest, _ := breakRuneSlice(value)
	sofar = append(sofar, first)

	ch, ok := n.children[first]
	if !ok {
		return false, sofar
	}

	if depth == 0 {
		if ch.isTerminated {
			return true, sofar
		}
	}

	if depth != 0 {
		depth--
	}

	return ch.isChildWithDepth(rest, depth, sofar)

}
