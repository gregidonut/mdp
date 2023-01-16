package main

import (
	"flag"
	"fmt"
	"os"
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
	return nil
}
