package symfony

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/lettland/cache-warmer/structs"
)

// FindFiles searches for files in the specified root directory and its subdirectories.
// It excludes the specified directories and handles vendor directories based on the vendorWatch flag.
// It returns a list of file paths and an error if any occurred.
func FindFiles(config structs.Config, root string, excludedDirs []string, vendorWatch bool, vendorList []string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip excluded directories (not individual files)
		for _, excludedDir := range excludedDirs {
			if d.IsDir() && strings.Contains(path, excludedDir) {
				return filepath.SkipDir
			}
		}

		// Handle the vendor directory based on vendorWatch
		if d.IsDir() {
			if strings.HasPrefix(path, filepath.Join(root, config.DirSymfonyVendor)) {
				if !vendorWatch {
					// If vendorWatch is false, skip the entire vendor directory
					return filepath.SkipDir
				}

				// If vendorWatch is true, check if the current path matches any of the vendorList entries
				if len(vendorList) > 0 {
					insideVendor := false
					for _, vendor := range vendorList {
						if strings.HasPrefix(path, filepath.Join(root, config.DirSymfonyVendor, vendor)) {
							insideVendor = true
							break
						}
					}
					if !insideVendor {
						return filepath.SkipDir
					}
				}
			}
		} else {
			// Add the file if it's not excluded (like .gitignore)
			if !strings.HasSuffix(path, ".gitignore") {
				files = append(files, path)
			}
		}

		return nil
	})

	return files, err
}

// GetWatchMap returns a map containing the files to watch and their corresponding last modified timestamps.
// It takes a `config` parameter of type `structs.Config` which holds the configuration values for the application.
// It calls the `GetFilesToWatch` function to retrieve the files to watch.
// For each file, it retrieves the file's stats using `os.Stat`, and adds the file path and last modified timestamp
// to the `watchMap`. If any error occurs during the process, it returns the error.
// It returns the `watchMap` containing the files to watch and their corresponding timestamps and nil error on success.
// Otherwise, it returns nil map and the encountered error.
// The `watchMap` can be used to compare with the existing files being watched to detect any changes.
//
// Note that `GetWatchMap` does not handle removing files from the `watchMap` when they are no longer being watched.
// This responsibility falls on the caller of this function.
func GetWatchMap(config structs.Config) (map[string]string, error) {
	watchMap := make(map[string]string)
	filesToWatch, err := GetFilesToWatch(config)

	if err != nil {
		return nil, err
	}

	for _, file := range filesToWatch {
		stats, err := os.Stat(file)
		if err != nil {
			return nil, fmt.Errorf("can't get stats for the \"%s\" file, check the project permissions or if a new file was created: %v", file, err)
		}
		watchMap[file] = stats.ModTime().String()
	}

	return watchMap, nil
}

func GetFilesToWatch(config structs.Config) ([]string, error) {
	var filesToWatch []string

	// Set up excluded directories
	var excludedDirs = config.DirsExclude
	if !config.VendorWatch {
		excludedDirs = append(excludedDirs, config.DirSymfonyVendor)
	}

	// Include general files like .env*
	envFiles, err := GetFilesFromPath(config, ".env*")
	if err != nil {
		return nil, err
	}
	filesToWatch = append(filesToWatch, envFiles...)

	indexFile, err := GetSingleFileFromPath(config, "public/index.php")
	if err != nil {
		return nil, err
	}
	filesToWatch = append(filesToWatch, indexFile)

	// Directories to watch
	symfonyDirs := map[string]string{
		config.DirSymfonyConfig:       config.DirSymfonyConfig,
		config.DirSymfonySrc:          config.DirSymfonySrc,
		config.DirSymfonyTemplates:    config.DirSymfonyTemplates,
		config.DirSymfonyTranslations: config.DirSymfonyTranslations,
		config.DirMigrations:          config.DirMigrations,
	}

	// Watch all files in the specified directories, regardless of their extensions
	for dir := range symfonyDirs {
		files, err := FindFiles(config, dir, excludedDirs, false, config.VendorList)
		if err != nil {
			return nil, err
		}
		filesToWatch = append(filesToWatch, files...)
	}

	// If VendorWatch is enabled, watch specific vendor directories
	if config.VendorWatch {
		for _, vendor := range config.VendorList {
			vendorPath := filepath.Join(config.DirSymfonyVendor, vendor)
			vendorFiles, err := FindFiles(config, vendorPath, excludedDirs, true, config.VendorList)
			if err != nil {
				return nil, err
			}
			filesToWatch = append(filesToWatch, vendorFiles...)
		}
	}

	return filesToWatch, nil
}

// GetFilesFromPath retrieves a list of files from the specified path based on the provided configuration.
// It takes a `config` parameter of type `structs.Config` which holds the configuration values.
// The function uses the filepath package to glob files based on the path and the provided glob pattern.
// If an error occurs during the globbing process, the function returns an error.
// Otherwise, it returns the list of files matching the pattern and nil error.
func GetFilesFromPath(config structs.Config, glob string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(config.DirSymfonyProject, glob))
	if err != nil {
		return nil, fmt.Errorf("error while globbing files: %v", err)
	}

	return files, nil
}

func GetSingleFileFromPath(config structs.Config, filePath string) (string, error) {
	fullPath := filepath.Join(config.DirSymfonyProject, filePath)

	// Check if the file exists
	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file not found: %s", fullPath)
		}
		return "", err
	}

	return fullPath, nil
}
