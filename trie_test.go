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
	"strings"
	"testing"
)

func TestTrieLoading(t *testing.T) {

	list := []string{"copy", "copper", "workflow", "workshop", "workbench", "work"}

	trie := New()

	if err := trie.Load(list); err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if trie.Count() != len(list) {
		t.Errorf("Expected %d, got %d", len(list), trie.Count())
	}

}

func TestTrieDelete(t *testing.T) {

	list := []string{"cop", "copy", "copper", "copperhead"}

	trie := New()

	cases := []struct {
		In     string
		Before bool
		After  bool
	}{
		{"cop", true, true},
		{"copy", true, true},
		{"copper", true, false},
		{"copperhead", true, true},
		{"CoPy", true, true},
		{"1copper", false, false},
	}

	if err := trie.Load(list); err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if trie.Count() != len(list) {
		t.Errorf("Expected %d, got %d", len(list), trie.Count())
	}

	for _, c := range cases {
		got := trie.Find(c.In)
		if c.Before != got {
			t.Errorf("For %s Expected %t, got %t", c.In, c.Before, got)
		}
	}

	if err := trie.Delete("copper"); err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if trie.Count() != len(list)-1 {
		t.Errorf("Expected %d, got %d", len(list)-1, trie.Count())
	}

	for _, c := range cases {
		got := trie.Find(c.In)
		if c.After != got {
			t.Errorf("For %s Expected %t, got %t", c.In, c.After, got)
		}
	}

}

func TestTrieLoadingEmpty(t *testing.T) {

	list := []string{}

	trie := New()

	if err := trie.Load(list); err != ErrTrieLoadEmpty {
		t.Errorf("Expected %v, got %v", ErrTrieLoadEmpty, err)
	}

}

func TestTrieLoadingFile(t *testing.T) {
	trie := New()

	if err := trie.LoadFile("dict.json"); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestTrieLoadingNoFile(t *testing.T) {
	trie := New()

	err := trie.LoadFile("dict_does_not_exists.json")
	if err != nil && strings.Index(err.Error(), "no such file") < 0 {
		t.Errorf("Expected 'no such file' error, got %v", err)
	}
}

func TestTrieLoadingBadFile(t *testing.T) {
	trie := New()

	err := trie.LoadFile("dict.bad.json")
	if err != nil && strings.Index(err.Error(), "cannot unmarshall") < 0 {
		t.Errorf("Expected 'unmarshalling json' error, got %v", err)
	}
}

func TestTrieFinding(t *testing.T) {

	list := []string{"copy", "copper", "workflow", "workshop", "workbench", "work", "a", "Apple", "appleseed"}

	cases := []struct {
		In  string
		Out bool
	}{
		{"copy", true},
		{"copper", true},
		{"copperhead", false},
		{"workflow", true},
		{"workshop", true},
		{"workbench", true},
		{"work", true},
		{"flow", false},
		{"failwork", false},
		{"space", false},
		{"CoPy", true},
		{"a", true},
		{"t", false},
		{"1copper", false},
		{"&5847234@#$@#$", false},
		{"cop", false},
		{"apple", true},
		{"tapple", false},
	}

	trie := New()

	if err := trie.Load(list); err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	for _, c := range cases {
		got := trie.Find(c.In)
		if c.Out != got {
			t.Errorf("For %s Expected %t, got %t", c.In, c.Out, got)
		}

	}

}

func TestTrieFindLoadOrderBug(t *testing.T) {

	list := []string{"workbench", "work"}

	cases := []struct {
		In  string
		Out bool
	}{
		{"work", true},
	}

	trie := New()

	if err := trie.Load(list); err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	for _, c := range cases {
		got := trie.Find(c.In)
		if c.Out != got {
			t.Errorf("For %s Expected %t, got %t", c.In, c.Out, got)
		}

	}

}

func TestTrieIsContained(t *testing.T) {

	list := []string{"copy", "copper", "workflow", "workshop", "workbench", "work", "a", "apple"}

	cases := []struct {
		In     string
		Report string
		Out    bool
	}{
		{"copy", "copy", true},
		{"copper", "copper", true},
		{"copperhead", "copper", true},
		{"workflow", "work", true},
		{"workshop", "work", true},
		{"workbench", "work", true},
		{"work", "work", true},
		{"flow", "", false},
		{"failwork", "work", true},
		{"space", "", false},
		{"CoPy", "copy", true},
		{"1copper", "copper", true},
		{"&5847234@#$@#$", "", false},
		{"&5847copper234@#$@#$", "copper", true},
		{"copper234@#$@#$", "copper", true},
		{"&5847copper", "copper", true},
		{"Drdfjflr9mg&Apple", "apple", true},
	}

	trie := New()

	if err := trie.Load(list); err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	for _, c := range cases {
		got, gotw := trie.IsContained(c.In, 3)
		if c.Out != got {
			t.Errorf("For %s Expected %t, got %t", c.In, c.Out, got)
		}

		if c.Report != gotw {
			t.Errorf("For %s Expected %s, got %s", c.In, c.Report, gotw)
		}

	}

}

func TestTrieIsContainedSubStringBug(t *testing.T) {

	list := []string{"a", "cope", "copper", "zzz"}

	cases := []struct {
		In     string
		Report string
		Out    bool
	}{
		{"copy", "", false},
		{"copper", "copper", true},
		{"copperhead", "copper", true},
		{"workflow", "", false},
		{"workshop", "", false},
		{"workbench", "", false},
		{"work", "", false},
		{"flow", "", false},
		{"failwork", "", false},
		{"space", "", false},
		{"CoPy", "", false},
		{"1copper", "copper", true},
		{"&5847234@#$@#$", "", false},
		{"&5847copper234@#$@#$", "copper", true},
		{"copper234@#$@#$", "copper", true},
		{"&5847copper", "copper", true},
		{"Drdfjflr9mg&Apple", "", false},
	}

	trie := New()

	if err := trie.Load(list); err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	for _, c := range cases {
		got, gotw := trie.IsContained(c.In, 3)
		if c.Out != got {
			t.Errorf("For %s Expected %t, got %t", c.In, c.Out, got)
		}

		if c.Report != gotw {
			t.Errorf("For %q Expected %q, got %q", c.In, c.Report, gotw)
		}

	}

}

func BenchmarkSearch(b *testing.B) {
	trie := New()

	if err := trie.LoadFile("dict.full.json"); err != nil {
		b.Errorf("Expected no error, got %v", err)
	}

	data, err := fileToStringSlice("dict.full.json")

	if err != nil {
		b.Fatalf("Error in reading in file for testing %v", err)
	}

	dictloop := 0

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		if dictloop >= len(data) {
			dictloop = 0
		}
		trie.Find(data[dictloop])
		dictloop++
	}
}

func BenchmarkContains(b *testing.B) {
	trie := New()

	if err := trie.LoadFile("dict.full.json"); err != nil {
		b.Errorf("Expected no error, got %v", err)
	}

	data, err := fileToStringSlice("dict.full.json")

	if err != nil {
		b.Fatalf("Error in reading in file for testing %v", err)
	}

	dictloop := 0

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		if dictloop >= len(data) {
			dictloop = 0
		}
		trie.IsContained(data[dictloop], 0)
		dictloop++
	}
}
