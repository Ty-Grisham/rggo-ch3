package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	// Defining a boolean flag -l to count lines instead of words
	lines := flag.Bool("l", false, "Count lines")
	bytes := flag.Bool("b", false, "Count Bytes")
	// Parsing the flags provided by the user
	flag.Parse()

	// Calling the count function to count the number of words (or lines)
	// recieved from the STDIN and printing it out
	fmt.Println(count(os.Stdin, *lines, *bytes))
}

// Count counts the number of words in an input io.reader object
func count(r io.Reader, countLines, countBytes bool) int {
	// A scanner is used to read text from a Reader (such as files)
	scanner := bufio.NewScanner(r)

	// If countLines & countBytes is false, define the scanner split-type to words
	if !countLines {
		scanner.Split(bufio.ScanWords)
	}

	switch {
	case countLines:
		scanner.Split(bufio.ScanLines)
	case countBytes:
		scanner.Split(bufio.ScanBytes)
	default:
		scanner.Split(bufio.ScanWords)
	}
	// Defining a counter
	wc := 0

	// For every word or line scanned, increment the counter
	for scanner.Scan() {
		wc++
	}

	// Return the total
	return wc
}
