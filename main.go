package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

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
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	flag.Parse()

	// If user did not provide input, show usage
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}
	if err := run(*filename, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run coordinates the execution of multiple functions
func run(filename string, out io.Writer, skipPreview bool) error {
	// Read all data from the input file and check for errors
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData := parseContent(input)

	// Create temporary file and check for errors
	temp, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}

	outName := temp.Name()

	fmt.Fprintln(out, outName)

	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	defer os.Remove(outName)

	return preview(outName)
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

// preview automatically opens an HTML file. preview takes the
// temporary file name as input and returns an error in case it can
// not open the file
func preview(fname string) error {
	cName := ""
	cParams := []string{}

	// Define executeable based on OS
	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("os not supported")
	}

	// Append file name to parameters slice
	cParams = append(cParams, fname)

	// Locate executable in PATH
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	// Open the file using default programs
	err = exec.Command(cPath, cParams...).Run()

	// Give the browser some time to open the file bfore deleting it
	time.Sleep(2 * time.Second)
	return err
}
