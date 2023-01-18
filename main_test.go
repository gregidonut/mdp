package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const (
	inputFile                   = "./testdata/test1.md"
	goldenFileDefault           = "./testdata/test1.md.html"
	templateFile1               = "./testdata/testTemplate1.html"
	goldenFileWithTemplateFile1 = "./testdata/test2.md.html"
	templateFile2               = "./testdata/testTemplate2.html"
	goldenFileWithTemplateFile2 = "./testdata/test3.md.html"
)

func Test_parseContent(t *testing.T) {
	input, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}

	got, err := parseContent(input, inputFile, "")
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

var binName = "mdp"

// TestMain executes the go build tool to create the binary of the cli
// then runs it with arguments to check if outputs are expected
// then cleans up the files after the function is completed.
func TestMain(m *testing.M) {
	fmt.Println("Building tool...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	build := exec.Command("go", "build", "-o", binName)
	err := build.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %q: %s", binName, err)
	}

	fmt.Println("Running tests...")
	result := m.Run()

	fmt.Println("Cleaning up...")
	os.Remove(binName)

	os.Exit(result)
}

func TestMDPCLI(t *testing.T) {
	tests := []struct {
		name           string
		goldenFileName string
		flags          []string
		specifyEnvVar  bool
	}{
		{
			name:           "WithoutEnvVar -s -file " + inputFile,
			flags:          []string{"-s", "-file", inputFile},
			goldenFileName: goldenFileDefault,
		},
		{
			name:           "WithoutEnvVar -s -file " + inputFile + " -t " + templateFile1,
			flags:          []string{"-s", "-file", inputFile, "-t", templateFile1},
			goldenFileName: goldenFileWithTemplateFile1,
		},
		{
			name:           "WithEnvVar -s -file " + inputFile,
			flags:          []string{"-s", "-file", inputFile},
			goldenFileName: goldenFileWithTemplateFile2,
			specifyEnvVar:  true,
		},
		{
			name:           "WithoutEnvVar -s -file " + inputFile + " -t " + templateFile1,
			flags:          []string{"-s", "-file", inputFile, "-t", templateFile1},
			goldenFileName: goldenFileWithTemplateFile1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.specifyEnvVar {
				err := os.Setenv("MDP_TEMPLATE", templateFile2)
				if err != nil {
					t.Fatal(err)
				}
			}

			dir, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}

			cmdPath := filepath.Join(dir, binName)
			cmd := exec.Command(cmdPath, tt.flags...)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Log(string(out))
				t.Fatal(err)
			}

			t.Logf("output: \n%#q\n", string(out))

			outputFileContents, err := os.ReadFile(strings.TrimSpace(string(out)))
			if err != nil {
				log.Fatal(err)
			}

			goldenFileContents, err := os.ReadFile(tt.goldenFileName)
			if err != nil {
				log.Fatal(err)
			}

			got := string(outputFileContents)
			want := string(goldenFileContents)

			if got != want {
				t.Errorf("got != want, got: \n%s\n golden: %q \n%s\n", got, tt.goldenFileName, want)
			}
		})
	}
}
