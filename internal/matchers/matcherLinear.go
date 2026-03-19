package matchers

import "strings"

type MatcherLinear struct {
	matchPath bool
}

func (m MatcherLinear) Match(dirs []string, needle string) []string {
	res := make([]string, 0, len(dirs))

	for _, hay := range dirs {
		if !isMatch(strings.ToLower(hay), strings.ToLower(needle), m.matchPath) {
			continue
		}

		res = append(res, hay)
	}

	if len(res) == 0 {
		return dirs
	}

	return res
}

func isMatch(haystack string, needle string, matchPath bool) bool {
	for i := len(haystack) - 1; i >= len(needle)-1; i-- {
		if !matchPath && haystack[i] == '/' {
			return false
		}

		j := 0
		for j < len(needle) && haystack[i-j] == needle[len(needle)-j-1] {
			j++
		}

		if j == len(needle) {
			return true
		}
	}

	return false
}
