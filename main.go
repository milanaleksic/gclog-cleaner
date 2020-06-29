package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
)

const (
	Reading int = iota
	PassThrough
)

var (
	beginMatcher *regexp.Regexp
	inputFile    *os.File
	exclusions   = []*regexp.Regexp{
		regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.*\[GC`),
	}
)

type exclusionsType []string

func (i *exclusionsType) String() string {
	return "exclusions (regex patterns)"
}

func (i *exclusionsType) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func init() {
	var beginPattern string
	var inputFileLocation string
	var exclusionsPatterns exclusionsType
	flag.StringVar(&beginPattern, "begin-pattern", `\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`, "pattern that should match beginning of all log lines")
	flag.StringVar(&inputFileLocation, "input-file", "", "which file to process (default - stdin)")
	flag.Var(&exclusionsPatterns, "exclusions", "exclusion patterns (default: only one - '\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.*\\[GC')")
	flag.Parse()

	beginMatcher = regexp.MustCompile(beginPattern)
	if inputFileLocation != "" {
		var err error
		inputFile, err = os.Open(inputFileLocation)
		if err != nil {
			log.Fatalf("Failed to open input file: %s, reason: %v", inputFileLocation, err)
		}
	}
	exclusions = make([]*regexp.Regexp, len(exclusionsPatterns))
	for i, pattern := range exclusionsPatterns {
		exclusions[i] = regexp.MustCompile(pattern)
	}
}

func main() {
	var s int
	var scanner *bufio.Scanner
	if inputFile != nil {
		scanner = bufio.NewScanner(inputFile)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
		_, _ = fmt.Fprintln(os.Stderr, "Reading from stdin")
	}
	for scanner.Scan() {
		line := scanner.Text()
		if filteredOut(line) {
			s = PassThrough
		} else if isLogLine(line) {
			s = Reading
		}
		if s == Reading {
			fmt.Println(line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func filteredOut(text string) bool {
	for _, exclusion := range exclusions {
		if len(exclusion.FindAllString(text, -1)) > 0 {
			return true
		}
	}
	return false
}

func isLogLine(text string) bool {
	return len(beginMatcher.FindAllString(text, -1)) > 0
}
