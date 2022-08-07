package main

import (
	"testing"
)

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
