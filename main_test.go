package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

const (
	inputFile  = "./testdata/test1.md"
	goldenFile = "./testdata/test1.md.html"
)

func Test_parseContent(t *testing.T) {
	input, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}

	got := parseContent(input)

	want, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(want, got) {
		t.Logf("golden:\n%s\n", want)
		t.Logf("result:\n%s\n", got)
		t.Error("Result content does not match golden file")
	}
}

func TestRun(t *testing.T) {
	mockStdOut := bytes.Buffer{}

	if err := run(inputFile, &mockStdOut, true); err != nil {
		t.Fatal(err)
	}

	resultFile := strings.TrimSpace(mockStdOut.String())

	want, err := os.ReadFile(resultFile)
	if err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(want, got) {
		t.Logf("golden:\n%s\n", want)
		t.Logf("result:\n%s\n", got)
		t.Error("Result content does not match golden file")
	}

	os.Remove(resultFile)
}
