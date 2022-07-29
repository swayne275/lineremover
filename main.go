package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	filePath := flag.String("file", "", "file to modify")
	keys := flag.String("keys", "", "keys to search for in lines - separate multiple keys with '|'")
	inPlace := flag.Bool("inplace", false, "edit the file (don't create a copy)")
	flag.Parse()

	if *filePath == "" || *keys == "" {
		helpAndExit()
	}
	keyPhrases := strings.Split(*keys, "|")

	log.Printf("Trimming file %q of lines with key phrases: %#v (in place: %t)\n", *filePath, keyPhrases, *inPlace)

	if err := cut(*filePath, keyPhrases, *inPlace); err != nil {
		log.Fatalf("failed to cut lines: %s", err)
	}
}

// pretty print tool instructions, then exit the program.
// use fmt package to avoid log prefixes in the message.
func helpAndExit() {
	fmt.Println(`Line remover tool help:

Supply a -file to tell the program which file to modify (relative paths work).

Supply the -keys to search for. If a line in the -file contains at least
one of these, it will be removed. Multiple keys may be separated by a '|'.

Optionally, set -inplace=true to perform the operation in-place (edit
the provided -file rather than creating a new one).`)

	os.Exit(1)
}

func cut(filePath string, keyPhrases []string, inplace bool) error {
	tempFilePath := getTempFilePath(filePath)

	err := cutLines(filePath, tempFilePath, keyPhrases)
	if err != nil {
		os.Remove(tempFilePath)
		return err
	}

	if inplace {
		return os.Rename(tempFilePath, filePath)
	}

	return nil
}

// generate a temp file that includes all of `filePath` minus lines matching a `keyPhrase`
func cutLines(filePath, tempFilePath string, keyPhrases []string) (retErr error) {
	sourceFile, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed opening source file: %s", err)
	}
	defer sourceFile.Close()

	outFile, err := os.OpenFile(tempFilePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("failed creating output file: %s", err)
	}
	defer func() {
		e := outFile.Close()
		if retErr == nil {
			retErr = e
		}
	}()
	bw := bufio.NewWriter(outFile)
	defer func() {
		e := bw.Flush()
		if retErr == nil {
			retErr = e
		}
	}()

	first := true
	scanner := bufio.NewScanner(sourceFile)
	for scanner.Scan() {
		line := scanner.Text()
		if !substrInLine(line, keyPhrases) {
			b := strings.Builder{}
			if !first {
				b.WriteString("\n")
			}
			b.WriteString(line)
			if _, err := bw.WriteString(b.String()); err != nil {
				log.Println("Error writing string to buffered output:", err)
				return err
			}

			first = false
		}
	}

	return nil
}

// TODO clean up code when all done

// check if any of `keyPhrases`` are in `line`
// TODO might be best to compile the keyphrases to a regex then check against that.
// i'll need to benchmark to see which is best
// TODO variadic keyphrases as ...string instead of []string?
func substrInLine(line string, keyPhrases []string) bool {
	for _, keyPhrase := range keyPhrases {
		if strings.Contains(line, keyPhrase) {
			return true
		}
	}

	return false
}

// TODO this should be reformed to make the regex, then search all lines
func substrInLineRegex(line string, keyPhrases []string) bool {
	re, err := buildRegex(keyPhrases)
	if err != nil {
		panic(fmt.Sprintf("TODO cleanup error compiling regex: %v", err))
	}

	return re.MatchString(line)
}

func buildRegex(keyPhrases []string) (*regexp.Regexp, error) {
	regexSB := strings.Builder{}
	for _, keyphrase := range keyPhrases {
		if regexSB.Len() > 0 {
			regexSB.WriteRune('|')
		}
		regexSB.WriteString(keyphrase)
	}

	return regexp.Compile(regexSB.String())
}

// generate the file path for the temporary file
func getTempFilePath(filePath string) string {
	return filePath + ".tmp"
}

// convenience function to get the final file path
func getOutputFilePath(filePath string, inPlace bool) string {
	if inPlace {
		return filePath
	}

	return getTempFilePath(filePath)
}