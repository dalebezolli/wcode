package matchers

import (
	"fmt"
	"regexp"
	"strings"
)

type MatcherRG struct {
}

func (m MatcherRG) Match(dirs []string, needle string) []string {
	matches := make([]string, 0, len(dirs))

	expr := fmt.Sprintf("/[^/]*%[1]s[^/]*$|/[^/]*%[1]s[^/]*/[^/]*$", strings.ReplaceAll(needle, " ", ".*"))
	r, err := regexp.Compile(expr)
	if err != nil {
		return dirs
	}

	for _, dir := range dirs {
		if r.Match([]byte(dir)) {
			matches = append(matches, dir)
		}
	}

	if len(matches) == 0 {
		return dirs
	}

	return matches
}
