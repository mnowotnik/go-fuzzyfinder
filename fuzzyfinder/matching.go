package fuzzyfinder

import (
	"sort"
	"strings"
	"unicode"
)

// Matched represents a result of FindAll.
type Matched struct {
	// Idx is the index of an item of the original slice which was used to
	// search matched strings.
	Idx int
	// Pos is the range of matched position.
	// [2]int represents an open interval of a position.
	Pos [2]int
	// score is the value that indicates how it similar to the input string.
	// The bigger score, the more similar it is.
	score int
}

func findAll(in string, slice []Item, opts ...Option) []Matched {
	var opt opt
	for _, o := range opts {
		o(&opt)
	}
	m := match(in, slice, opt)
	sort.Slice(m, func(i, j int) bool {
		return m[i].score > m[j].score
	})
	return m
}

// match iterates each string of slice for check whether it is matched to the input string.
func match(input string, slice []Item, opt opt) (res []Matched) {
	if opt.mode == ModeSmart {
		// Find an upper-case rune
		n := strings.IndexFunc(input, unicode.IsUpper)
		if n == -1 {
			opt.mode = ModeCaseInsensitive
			input = strings.ToLower(input)
		} else {
			opt.mode = ModeCaseSensitive
		}
	}

	in := []rune(input)
	for idxOfSlice, item := range slice {
		var idx int
		var s string
		if item.View != "" {
			s = item.View
		} else {
			s = item.Value
		}
		if opt.mode == ModeCaseInsensitive {
			s = strings.ToLower(s)
		}
	LINE_MATCHING:
		for _, r := range []rune(s) {
			if r == in[idx] {
				idx++
				if idx == len(in) {
					score, pos := Calculate(s, input)
					res = append(res, Matched{
						Idx:   idxOfSlice,
						Pos:   pos,
						score: score,
					})
					break LINE_MATCHING
				}
			}
		}
	}
	return
}
