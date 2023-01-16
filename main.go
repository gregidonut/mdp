package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"os"
	"path/filepath"
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
	filename := flag.String("file", "", "Markdown file preview")

	flag.Parse()

	// Put flag in a container that we can check later if the flag is used
	usedFlags := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		usedFlags[f.Name] = true
	})

	switch {
	case usedFlags["file"]:
		if *filename == "" {
			flag.Usage()
			os.Exit(1)
		}
	}

	if err := run(*filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(fileName string) error {
	// Read all the data from the input file and check for errors
	input, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	htmlData := parseContent(input)
	outName := fmt.Sprintf("%s.html", filepath.Base(fileName))
	fmt.Println(outName)

	return saveHTML(outName, htmlData)
}

func parseContent(input []byte) []byte {
	// Parse the markdown file through blackfriday and bluemonday
	// to generate a valid and safe HTML
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	// generate html
	var buffer bytes.Buffer
	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)
	return buffer.Bytes()
}

func saveHTML(outFName string, data []byte) error {
	return os.WriteFile(outFName, data, 0644)
}
