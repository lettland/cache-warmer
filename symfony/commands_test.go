package symfony

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/lettland/cache-warmer/structs"
)

func TestCheckSymfonyConsole(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tmpDir)

	var testCases = map[string]struct {
		consolePath   string
		expectedError string
		createConsole bool
	}{
		"ConsoleExists": {
			consolePath:   "bin/symfony-console",
			expectedError: "",
			createConsole: true,
		},
		"ConsoleMissing": {
			consolePath:   "undefined/symfony-console",
			expectedError: fmt.Sprintf("symfony console not found at %s", filepath.Join(tmpDir, "undefined", "symfony-console")),
			createConsole: false,
		},
	}

	for key, tc := range testCases {
		t.Run(key, func(t *testing.T) {
			if tc.createConsole {
				err := os.MkdirAll(filepath.Join(tmpDir, filepath.Dir(tc.consolePath)), os.ModePerm)
				if err != nil {
					t.Fatal(err)
				}
				_, err = os.Create(filepath.Join(tmpDir, tc.consolePath))
				if err != nil {
					t.Fatal(err)
				}
			}

			config := structs.Config{
				DirSymfonyProject:  tmpDir,
				SymfonyConsolePath: tc.consolePath,
			}

			err := CheckSymfonyConsole(config)
			if err != nil && err.Error() != tc.expectedError {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestGetSymfonyProjectDir(t *testing.T) {
	execDir, _ := os.Getwd()
	osArgs := os.Args
	defer func() { os.Args = osArgs }()

	testCases := map[string]struct {
		args          []string
		expectedDir   string
		expectedError string
		preFunc       func()
	}{
		"ExistingAbsolutePath": {
			args:        []string{"", "/tmp"},
			expectedDir: "/tmp",
		},
		"ExistingRelativePath": {
			args:        []string{"", "symfony"},
			expectedDir: filepath.Join(execDir, "symfony"),
			preFunc: func() {
				err := EnsureDir("symfony")
				if err != nil {
					t.Fatal(err)
				}
			},
		},
		"MissingPath": {
			args:          []string{"", "missingDir"},
			expectedError: fmt.Sprintf("stat %s: no such file or directory", filepath.Join(execDir, "missingDir")),
		},
		"NoArguments": {
			args:          []string{""},
			expectedError: "no path provided",
		},
	}

	for testName, tc := range testCases {
		os.Args = tc.args
		if tc.preFunc != nil {
			tc.preFunc()
		}

		result, err := GetSymfonyProjectDir()
		if err != nil && err.Error() != tc.expectedError {
			t.Errorf("%s: unexpected error, expected %v, but got: %v", testName, tc.expectedError, err)
		}
		if tc.expectedDir != result {
			t.Errorf("%s: unexpected directory return, expected %v, but got: %v", testName, tc.expectedDir, result)
		}
	}
}
