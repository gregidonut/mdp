package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"html/template"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var (
	defaultTemplate = `<!DOCTYPE html>
<html>
    <head>
        <meta http-equiv="content-type" content="text/html; charset=utf-8">
        <title>{{ .Title }}</title>
    </head>
    <body>
    <header><h1>{{ .FileName }}</h1></header>
{{ .Body }}
    </body>
</html>
`
)

// content represents the HTML content to add into the template
type content struct {
	Title    string
	FileName template.HTML
	Body     template.HTML
}

func main() {

	// Parse flags
	flag.Bool("s", false, "Skip auto-preview")
	flag.Bool("stdin", false, "input from stdin")
	filename := flag.String("file", "", "Markdown file preview")
	tFName := flag.String("t", "", "Alternate template file name")

	flag.Parse()

	// Put flag in a container that we can check later if the flag is used
	usedFlags := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		usedFlags[f.Name] = true
	})

	if usedFlags["file"] {
		if *filename == "" {
			flag.Usage()
			os.Exit(1)
		}
	}

	// Check if the user defined the ENV VAR for custom template file
	if os.Getenv("MDP_TEMPLATE") != "" && !usedFlags["t"] {
		tfNameDeref := os.Getenv("MDP_TEMPLATE")
		// this needs to be a pointer to be compatible with flag.String() output
		tFName = &tfNameDeref
	}

	if err := run(*filename, *tFName, os.Stdout, usedFlags["s"], usedFlags["stdin"]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(fileName, tFName string, out io.Writer, skipPreview, inputFromSTDN bool) error {
	var input []byte

	if inputFromSTDN {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			scannerBytes := append(scanner.Bytes(), '\n')
			input = append(input, scannerBytes...)
		}
		fileName = "stdin"
	} else {
		var err error
		input, err = os.ReadFile(fileName)
		if err != nil {
			return err
		}
	}

	// Read all the data from the input file and check for errors
	htmlData, err := parseContent(input, fileName, tFName)
	if err != nil {
		return err
	}

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

func parseContent(input []byte, inputFName, tFName string) ([]byte, error) {
	// Parse the markdown file through blackfriday and bluemonday
	// to generate a valid and safe HTML
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	// Parse the contents of the defaultTemplate const into a new Template
	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}

	// If user provided alternate template file, replace template
	if tFName != "" {
		t, err = template.ParseFiles(tFName)
		if err != nil {
			return nil, err
		}
	}

	c := content{
		Title:    "Markdown Preview Tool",
		FileName: template.HTML(inputFName),
		Body:     template.HTML(body),
	}

	// generate html
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func saveHTML(outFName string, data []byte) error {
	return os.WriteFile(outFName, data, 0644)
}

func preview(fName string) error {
	cName := ""
	cParams := make([]string, 0)

	//Define executable based on OS
	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}

	// Append filename to parameters slice
	cParams = append(cParams, fName)

	// Locate executable in PATH
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	// Open the file using default program
	err = exec.Command(cPath, cParams...).Run()

	// Give the browser some time to open the file before deleting it
	time.Sleep(2 * time.Second)

	return err
}
