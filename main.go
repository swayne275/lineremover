package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	filePath := flag.String("file", "", "file to modify")
	keys := flag.String("keys", "", "keys to search for in lines - separate multiple keys with '|'")
	inPlace := flag.Bool("inplace", false, "edit the file (don't create a copy)")
	flag.Parse()

	keyPhrases := splitKeys(*keys)
	fmt.Printf("Trimming file '%s' of lines with key phrases: %#v\n", *filePath, keyPhrases)
	fmt.Printf("in place: %t\n", *inPlace)

	if err := cut(*filePath, keyPhrases, *inPlace); err != nil {
		log.Fatalf("failed to cut lines: %s", err)
	}
}

func splitKeys(keys string) []string {
	if keys == "" {
		return []string{}
	}
	return strings.Split(keys, "|")
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

// check if any of <keyPhrases> are in <line>
func substrInLine(line string, keyPhrases []string) bool {
	for _, keyPhrase := range keyPhrases {
		if strings.Contains(line, keyPhrase) {
			return true
		}
	}

	return false
}

// generate a temp file that includes all of <filePath> minus lines matching a <keyPhrase>
func cutLines(filePath, tempFilePath string, keyPhrases []string) error {
	sourceFile, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed opening source file: %s", err)
	}
	defer sourceFile.Close()

	outFile, err := os.Create(tempFilePath)
	if err != nil {
		log.Fatalf("failed opening output file: %s", err)
	}

	scanner := bufio.NewScanner(sourceFile)
	for scanner.Scan() {
		line := scanner.Text()
		if !substrInLine(line, keyPhrases) {
			// TODO use bufio?
			_, err := outFile.WriteString(line + "\n")
			if err != nil {
				outFile.Close()
				return err
			}
		}
	}

	return outFile.Close()
}

func cut(filePath string, keyPhrases []string, inplace bool) error {
	tempFilePath := getTempFilePath(filePath)
	err := cutLines(filePath, tempFilePath, keyPhrases)
	if err != nil {
		os.Remove(tempFilePath)
		return err
	}

	if inplace {
		// overwrite og file with temp file, return error if needed
		return nil
	}

	return nil
}
