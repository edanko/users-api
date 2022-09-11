package tr

import (
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const (
	shortStringLevenshteinLimit  uint8 = 0
	mediumStringLevenshteinLimit uint8 = 1
	longStringLevenshteinLimit   uint8 = 2

	shortStringThreshold  uint8 = 0
	mediumStringThreshold uint8 = 3
	longStringThreshold   uint8 = 5
)

type Entry[T any] struct {
	Name  string
	Entry T
}

// Trie is a data structure for storing common prefixes to strings for efficient comparison
// and retrieval.
type Trie[T any] struct {
	root                             *node[Entry[T]]
	fuzzy, normalised, caseSensitive bool
	levenshteinScheme                map[uint8]uint8
	levenshteinIntervals             []uint8
	// originalDict is a mapping of normalised to original string.
	originalDict map[string][]Entry[T]

	transformer transform.Transformer
}

// node is a node in a Trie which contains a map of runes to more node pointers
// if word is non-empty, this indicates that the node defines the end of a word
type node[T any] struct {
	children map[rune]*node[Entry[T]]
	word     string
}

type score struct {
	levenshtein uint8
	fuzzy       bool
}

// New creates a new empty trie. By default, fuzzy search is on and string normalisation is on.
// The default levenshtein scheme is on, where search strings of len 1-2 characters allow no
// distance, search strings of length 3-4 allow a levenshtein distance of 1, and search strings
// of length 5 or more runes allow a levenshtein distance of two.
func New[T any]() *Trie[T] {
	t := new(Trie[T])
	t.root = new(node[Entry[T]])
	t.root.children = make(map[rune]*node[Entry[T]])
	t.originalDict = make(map[string][]Entry[T])
	t.transformer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	t.WithFuzzy()
	t.WithNormalisation()
	t.DefaultLevenshtein()
	t.CaseInsensitive()
	return t
}

// WithFuzzy sets the Trie to use fuzzy matching on search.
func (t *Trie[T]) WithFuzzy() *Trie[T] {
	t.fuzzy = true
	return t
}

// WithoutFuzzy sets the Trie not to use fuzzy matching on search.
func (t *Trie[T]) WithoutFuzzy() *Trie[T] {
	t.fuzzy = false
	return t
}

// WithNormalisation sets the Trie to use normalisation on search.
// For example, Jurg will find Jürgen, Jürg will find Jurgen.
func (t *Trie[T]) WithNormalisation() *Trie[T] {
	t.normalised = true
	return t
}

// WithoutNormalisation sets the Trie not to use normalisation on search.
// for example Jurg won't find Jürgen, Jürg won't find Jurgen.
func (t *Trie[T]) WithoutNormalisation() *Trie[T] {
	t.normalised = false
	return t
}

// CaseSensitive sets the Trie to use case-sensitive search.
func (t *Trie[T]) CaseSensitive() *Trie[T] {
	t.caseSensitive = true
	return t
}

// CaseInsensitive sets the Trie to use case-insensitive search.
func (t *Trie[T]) CaseInsensitive() *Trie[T] {
	t.caseSensitive = false
	return t
}

// WithoutLevenshtein sets the Trie not to allow any levenshtein distance between
// the search string and any matches.
func (t *Trie[T]) WithoutLevenshtein() *Trie[T] {
	t.levenshteinScheme = map[uint8]uint8{0: 0}
	t.levenshteinIntervals = []uint8{0}
	return t
}

// DefaultLevenshtein sets the trie to use the default levenshtein scheme.
func (t *Trie[T]) DefaultLevenshtein() *Trie[T] {
	t.levenshteinScheme = map[uint8]uint8{
		shortStringThreshold:  shortStringLevenshteinLimit,
		mediumStringThreshold: mediumStringLevenshteinLimit,
		longStringThreshold:   longStringLevenshteinLimit}
	t.levenshteinIntervals = []uint8{longStringThreshold, mediumStringThreshold, longStringThreshold}
	return t
}

// CustomLevenshtein sets up a custom levenshtein scheme.
// WARNING, this function will panic if the scheme is invalid.
// A valid scheme is a series of pairs of search string length -> levenshtein distance.
// There must be one entry with zero as search string length.
func (t *Trie[T]) CustomLevenshtein(scheme map[uint8]uint8) *Trie[T] {
	_, ok := scheme[0]
	if !ok {
		panic("invalid levenshtein scheme for GAT")
	}
	t.levenshteinIntervals = make([]uint8, 0, len(scheme))
	for key := range scheme {
		t.levenshteinIntervals = append(t.levenshteinIntervals, key)
	}
	sort.Slice(t.levenshteinIntervals, func(i, j int) bool {
		return t.levenshteinIntervals[i] > t.levenshteinIntervals[j]
	})
	t.levenshteinScheme = scheme
	return t
}

// Insert inserts strings into the Trie
func (t *Trie[T]) Insert(entries ...Entry[T]) {
	for _, entry := range entries {
		switch {
		case t.normalised && t.caseSensitive:
			normal, _, err := transform.String(t.transformer, entry.Name)
			if err != nil {
				continue
			}
			t.originalDict[normal] = append(t.originalDict[normal], entry)
			entry.Name = normal
		case t.normalised && !t.caseSensitive:
			normal, _, err := transform.String(t.transformer, entry.Name)
			if err != nil {
				continue
			}
			normal = strings.ToLower(normal)
			t.originalDict[normal] = append(t.originalDict[normal], entry)
			entry.Name = normal
		case !t.normalised && !t.caseSensitive:
			lower := strings.ToLower(entry.Name)
			t.originalDict[lower] = append(t.originalDict[lower], entry)
			entry.Name = lower
		}
		currentNode := t.root
		for index, character := range entry.Name {
			child, ok := currentNode.children[character]
			if !ok {
				child = make(node[Entry[T]])
				child.children = make(map[rune]*node[Entry[T]])
				if index == len(entry.Name)-len(string(character)) {
					child.word = entry.Name
				}

				currentNode.children[character] = child
			}
			currentNode = child
		}
	}
}

// SearchAll is just like Search, but without a limit.
func (t *Trie[T]) SearchAll(search string) []T {
	return t.Search(search, 0)
}

// Search will return all complete words in the trie that have the search string as a prefix,
// taking into account the Trie's settings for normalisation, fuzzy matching and levenshtein distance scheme.
func (t *Trie[T]) Search(search string, limit int) []T {
	if len(search) == 0 {
		return []T{}
	}
	if t.normalised {
		var err error
		search, _, err = transform.String(t.transformer, search)
		if err != nil {
			return []T{}
		}
	}
	if !t.caseSensitive {
		search = strings.ToLower(search)
	}
	maxDistance := t.maxDistance(search)
	// start the recursive function
	collection := make(map[string]score)
	t.collect(collection, search, t.root, 0, maxDistance, limit, t.fuzzy, false)
	hits := make([]string, 0, len(collection))
	for key := range collection {
		hits = append(hits, key)
	}
	sort.Slice(hits, func(i, j int) bool {
		switch {
		case collection[hits[i]].levenshtein != collection[hits[j]].levenshtein:
			return collection[hits[i]].levenshtein < collection[hits[j]].levenshtein
		case collection[hits[i]].fuzzy && !collection[hits[j]].fuzzy:
			return false
		case !collection[hits[i]].fuzzy && collection[hits[j]].fuzzy:
			return true
		default:
			return hits[i] < hits[j]
		}
	})
	if len(hits) >= limit && limit != 0 {
		return hits[:limit]
	}
	if !t.normalised && t.caseSensitive {
		return hits
	}
	originals := make([]T, 0, len(hits)*2)
	for _, hit := range hits {
		originals = append(originals, t.originalDict[hit.Name]...)
	}
	return originals
}

// collect is a recursive function that traverses the Trie and inserts words from Word-final nodes which match the search
// text in the map collection. It handles substitution, insertion and deletion to the levenshtein distance limit and also
// allows fuzzy search.
func (t *Trie[T]) collect(collection map[string]score, word string, node *node[Entry[T]], distance, maxDistance uint8, limit int, fuzzyAllowed, fuzzyUsed bool) {
	if len(word) == 0 {
		if node.word != "" {
			previousScore, ok := collection[node.word]
			if !ok || distance < previousScore.levenshtein ||
				(distance == previousScore.levenshtein && previousScore.fuzzy && !fuzzyUsed) {
				collection[node.word] = score{levenshtein: distance, fuzzy: fuzzyUsed}
			}
			node.collectAllDescendentWords(collection, distance, fuzzyUsed)
			return
		}
		node.collectAllDescendentWords(collection, distance, fuzzyUsed)
	}
	character, size := utf8.DecodeRuneInString(word)
	subword := word[size:]
	// special rune for string collisions
	if character == '*' {
		t.collect(collection, subword, node, distance, maxDistance, limit, false, fuzzyUsed)
	}

	if next := node.children[character]; next != nil {
		t.collect(collection, subword, next, distance, maxDistance, limit, false, fuzzyUsed)
	}

	if distance < maxDistance {
		distance++

		for character, next := range node.children {
			// Substition
			t.collect(collection, string(character)+subword, node, distance, maxDistance, limit, false, fuzzyUsed)
			// Insertion
			t.collect(collection, string(character)+word, node, distance, maxDistance, limit, false, fuzzyUsed)
			// Fuzzy
			if fuzzyAllowed {
				t.collect(collection, word, next, distance-1, maxDistance, limit, true, true)
			}
		}
		// Deletion
		t.collect(collection, subword, node, distance, maxDistance, limit, false, false)
	} else if distance == 0 {
		for _, next := range node.children {
			// Fuzzy without levenshtein
			if fuzzyAllowed {
				t.collect(collection, word, next, distance, maxDistance, limit, true, true)
			}
		}
	}
}

// collectAllDescendentWords returns the words from all nodes that are descedent of the current node.
func (n *node[T]) collectAllDescendentWords(collection map[string]score, distance uint8, fuzzyUsed bool) {
	for _, node := range n.children {
		if node.word != "" {
			previousScore, ok := collection[node.word]

			if !ok || distance < previousScore.levenshtein ||
				(distance == previousScore.levenshtein && previousScore.fuzzy && !fuzzyUsed) {
				collection[node.word] = score{levenshtein: distance, fuzzy: fuzzyUsed}
			}
		}
		node.collectAllDescendentWords(collection, distance, fuzzyUsed)
	}
}

// maxDistance determines the maximum levenshein distance based on the levenshtein scheme
// and search string length.
func (t *Trie[T]) maxDistance(search string) (maxDistance uint8) {
	runes := []rune(search)
	for _, limit := range t.levenshteinIntervals {
		if len(runes) >= int(limit) {
			maxDistance = t.levenshteinScheme[limit]
			return
		}
	}
	return
}
