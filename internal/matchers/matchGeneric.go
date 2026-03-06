package matchers

type Matcher interface {
	Match(haystack []string, needle string) []string
}
