package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	bluemonday "github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	header = `<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="content-type" content="text/html; charset=utf-8">
    <title>Markdown Preview Tool</title>
  </head>
  <body>
`
	footer = `
  </body>
</html>
`
)

func main() {
	// Parse flags
	filename := flag.String("file", "", "Markdown file to preview")
	flag.Parse()

	// If user did not provide input, show usage
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}
	if err := run(*filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run coordinates the execution of multiple functions
func run(filename string) error {
	// Read all data from the input file and check for errors
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData := parseContent(input)

	outName := fmt.Sprintf("%s.html", filepath.Base(filename))
	fmt.Println(outName)

	return saveHTML(outName, htmlData)
}

// parseContent recieves a slice of bytes representing the content of a Markdown
// file and returns another slice of bytes with the converted content as HTML
func parseContent(input []byte) []byte {
	// Parse the Markdown file through blackfriday and bluemonday to
	// generate a safe and valid HTML
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	// Create a buffer of bytes to write to file
	var buffer bytes.Buffer

	// Write HTML to bytes buffer
	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)

	return buffer.Bytes()
}

// saveHTML recieves the entire HTML content stored in a bytes buffer and saves it
// to a file specified by outName
func saveHTML(outName string, data []byte) error {
	// Write the bytes to the file
	return os.WriteFile(outName, data, 0644)
}
