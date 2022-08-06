package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	cfg := getUserInput()

	log.Printf("Trimming file %q of lines with key phrases: %v (in place: %t)\n", cfg.inputPath, cfg.keys, cfg.inplace)

	if err := transformInput(cfg); err != nil {
		log.Fatalf("failed to cut lines: %s", err)
	}
}

type config struct {
	inputPath string
	keys      []string
	inplace   bool
}

func getUserInput() *config {
	inputPath := flag.String("file", "", "file to modify")
	keysRaw := flag.String("keys", "", "keys to search for in lines - separate multiple keys with '|'")
	inplace := flag.Bool("inplace", false, "edit the file (don't create a copy)")
	flag.Parse()

	if *inputPath == "" || *keysRaw == "" {
		helpAndExit()
	}
	keys := strings.Split(*keysRaw, "|")

	return &config{
		inputPath: *inputPath,
		keys:      keys,
		inplace:   *inplace,
	}
}

func transformInput(cfg *config) error {
	tempDstPath := cfg.inputPath + ".tmp"

	// clear any pre-existing output file
	os.Remove(tempDstPath)

	if err := transformInputImpl(cfg.inputPath, tempDstPath, cfg.keys); err != nil {
		// clean up after ourselves if there was an error
		os.Remove(tempDstPath)
		return err
	}

	if cfg.inplace {
		return os.Rename(tempDstPath, cfg.inputPath)
	}

	return nil
}

// generate a temp file that includes all of `inputPath` except for anything matching `keys`
func transformInputImpl(inputPath, tempDstPath string, keys []string) (retErr error) {
	sourceFile, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("failed opening source file: %v", err)
	}
	defer sourceFile.Close()

	outFile, err := os.OpenFile(tempDstPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("failed creating output file: %v", err)
	}
	defer func() {
		if e := outFile.Close(); retErr == nil {
			retErr = e
		}
	}()

	bw := bufio.NewWriter(outFile)
	defer func() {
		if e := bw.Flush(); retErr == nil {
			retErr = e
		}
	}()

	first := true
	scanner := bufio.NewScanner(sourceFile)
	for scanner.Scan() {
		line := scanner.Text()
		if !substrInLine(line, keys) {
			b := strings.Builder{}
			if !first {
				b.WriteString("\n")
			}
			b.WriteString(line)
			if _, err := bw.WriteString(b.String()); err != nil {
				return err
			}

			first = false
		}
	}

	return nil
}

// check if any of `keys` are in `line`
func substrInLine(line string, keys []string) bool {
	for _, key := range keys {
		if strings.Contains(line, key) {
			return true
		}
	}

	return false
}

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

// pretty print tool instructions, then exit the program.
// use fmt package to avoid log prefixes in the message.
func helpAndExit() {
	log.Println(`Line remover tool help:

Supply a -file to tell the program which file to modify (relative paths work).

Supply the -keys to search for. If a line in the -file contains at least
one of these, it will be removed. Multiple keys may be separated by a '|'.

Optionally, set -inplace=true to perform the operation in-place (edit
the provided -file rather than creating a new one).`)

	os.Exit(1)
}
