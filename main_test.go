package main

import (
	"regexp"
	"strings"
	"testing"
)

// note: regex is about 20x slower than simple substring search.
// decide if I want to have regex support or not
func buildRegex(keys []string) (*regexp.Regexp, error) {
	regexSB := strings.Builder{}
	for _, key := range keys {
		if regexSB.Len() > 0 {
			regexSB.WriteRune('|')
		}
		regexSB.WriteString(key)
	}

	return regexp.Compile(regexSB.String())
}

func BenchmarkRegex(b *testing.B) {
	re, err := buildRegex([]string{"hello", "tony", "world"})
	if err != nil {
		b.Fatalf("TODO cleanup error compiling regex: %v", err)
	}
	testStr := "hello world"

	for n := 0; n < b.N; n++ {
		re.MatchString(testStr)
	}
}

func BenchmarkStringSearch(b *testing.B) {
	testStr := "hello world"
	testPhrases := []string{"hello", "tony", "world"}
	for n := 0; n < b.N; n++ {
		substrInLine(testStr, testPhrases)
	}
}
