package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestParseCommaSeparated(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{name: "Empty string", input: "", want: []string{}},
		{name: "Single word", input: "test", want: []string{"test"}},
		{name: "Two words", input: "test,example", want: []string{"test", "example"}},
		{name: "Multiple words", input: "Universe,Earth,Mars,Jupiter,Saturn", want: []string{"Universe", "Earth", "Mars", "Jupiter", "Saturn"}},
		{name: "Leading and trailing spaces", input: " test , example ", want: []string{" test ", " example "}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseCommaSeparated(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCommaSeparated() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestGenerateSeparator(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  string
	}{
		{name: "Zero length", input: 0, want: ""},
		{name: "Length of 5", input: 5, want: "—————"},
		{name: "Length of 1", input: 1, want: "—"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateSeparator(tt.input); got != tt.want {
				t.Errorf("GenerateSeparator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateVersionLink(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{name: "Semantic version", version: "1.2.3", want: fmt.Sprintf("%s/releases/tag/%s", repository, "1.2.3")},
		{name: "Nightly version", version: "nightly", want: fmt.Sprintf("%s/releases/tag/%s", repository, "nightly")},
		{name: "Branch version", version: "master", want: fmt.Sprintf("%s/commit/%s", repository, "master")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateVersionLink(tt.version); got != tt.want {
				t.Errorf("GenerateVersionLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSemanticVersion(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "Semantic version", input: "1.2.3", want: true},
		{name: "Semantic version with prefix v", input: "v1.2.3", want: true},
		{name: "Semantic version with suffix", input: "1.2.3-alpha", want: true},
		{name: "Non-semantic version", input: "1.2", want: false},
		{name: "Non-semantic version with letters", input: "1.2.3a", want: false},
		{name: "Empty string", input: "", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSemanticVersion(tt.input); got != tt.want {
				t.Errorf("IsSemanticVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrintError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "Valid error",
			err:  errors.New("test error"),
			want: "/!\\ test error /!\\",
		},
		{
			name: "Nil error",
			err:  nil,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Call PrintError function
			PrintError(tt.err)

			// Restore output and read captured output
			w.Close()
			out, _ := ioutil.ReadAll(r)
			os.Stdout = rescueStdout

			// Clean up captured output
			got := strings.Trim(string(out), "\n")

			if got != tt.want {
				t.Errorf("PrintError() = %v, want %v", got, tt.want)
			}
		})
	}
}
