package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

const (
	inputFile                   = "./testdata/test1.md"
	goldenFileDefault           = "./testdata/test1.md.html"
	templateFile1               = "./testdata/testTemplate1.html"
	goldenFileWithTemplateFile1 = "./testdata/test2.md.html"
)

func Test_parseContent(t *testing.T) {
	input, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}

	got, err := parseContent(input, "")
	if err != nil {
		t.Fatal(err)
	}

	want, err := os.ReadFile(goldenFileDefault)
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
	type testInputs struct {
		name       string
		tFName     string
		goldenFile string
	}

	tests := []testInputs{
		{
			name:       "WithoutTFName",
			tFName:     "",
			goldenFile: goldenFileDefault,
		},
		{
			name:       "WithTFName",
			tFName:     templateFile1,
			goldenFile: goldenFileWithTemplateFile1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockStdOut := bytes.Buffer{}

			if err := run(inputFile, tt.tFName, &mockStdOut, true); err != nil {
				t.Fatal(err)
			}

			resultFile := strings.TrimSpace(mockStdOut.String())

			want, err := os.ReadFile(resultFile)
			if err != nil {
				t.Fatal(err)
			}

			got, err := os.ReadFile(tt.goldenFile)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(want, got) {
				t.Logf("golden:\n%s\n", want)
				t.Logf("result:\n%s\n", got)
				t.Error("Result content does not match golden file")
			}

			//time.Sleep(tt.delay)

			os.Remove(resultFile)
		})
	}
}
