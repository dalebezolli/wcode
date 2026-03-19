package matchers

type Matcher interface {
	Match(haystack []string, needle string) []string
}

type MatcherType string

const (
	MatcherTypeLinear MatcherType = "linear"
	MatcherTypeRegex  MatcherType = "regex"
)

func NewMatcher(matcherType MatcherType) Matcher {
	switch matcherType {
	case MatcherTypeLinear:
		return MatcherLinear{}

	case MatcherTypeRegex:
		return MatcherRG{}
	}

	return MatcherLinear{}
}
