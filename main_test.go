package main

import (
	"regexp"
	"testing"
)

var inputText = []string{
	"hello world",
	"hello big world",
	"hello bright world",
	"hello small big world",
}

func BenchmarkSubstringHello(b *testing.B) {
	cfg := &config{
		keys: []string{"hello"},
	}
	for n := 0; n < b.N; n++ {
		for _, line := range inputText {
			cfg.lineMatches(line)
		}
	}
}

func BenchmarkRegexHello(b *testing.B) {
	cfg := &config{
		pattern: regexp.MustCompile("hello"),
	}
	for n := 0; n < b.N; n++ {
		for _, line := range inputText {
			cfg.lineMatches(line)
		}
	}
}

func BenchmarkSubstringBig(b *testing.B) {
	cfg := &config{
		keys: []string{"big"},
	}
	for n := 0; n < b.N; n++ {
		for _, line := range inputText {
			cfg.lineMatches(line)
		}
	}
}

func BenchmarkRegexBig(b *testing.B) {
	cfg := &config{
		pattern: regexp.MustCompile("big"),
	}
	for n := 0; n < b.N; n++ {
		for _, line := range inputText {
			cfg.lineMatches(line)
		}
	}
}

func BenchmarkSubstringMultiple(b *testing.B) {
	cfg := &config{
		keys: []string{"big|brig|bight|bright"},
	}
	for n := 0; n < b.N; n++ {
		for _, line := range inputText {
			cfg.lineMatches(line)
		}
	}
}

func BenchmarkRegexMultiple(b *testing.B) {
	cfg := &config{
		pattern: regexp.MustCompile(".*b([r]?)ig([ht]?).*"),
	}
	for n := 0; n < b.N; n++ {
		for _, line := range inputText {
			cfg.lineMatches(line)
		}
	}
}
