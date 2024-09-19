package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"

	"github.com/lettland/cache-warmer/structs"
	"github.com/lettland/cache-warmer/symfony"
)

var version = "nightly"

const (
	repository = "https://github.com/lettland/cache-warmer"
)

// MainLoop continuously monitors for file changes and performs cache warming if an update is detected.
// It takes a `config` parameter of type `structs.Config` which holds the configuration values for the application.
// It also takes a `filesToWatch` parameter of type map[string]string that represents the files to watch for changes.
// The function checks for updated files using `symfony.GetWatchMap` and compares it with the existing `filesToWatch` map.
// If there are any differences, it starts cache warming by calling `symfony.CacheWarmup`. It measures the time taken
// to warm up the cache and prints the result. The updated `filesToWatch` map is then assigned to `filesToWatch`.
// If there are no differences, the function sleeps for a specified duration defined in the `config` parameter.
func MainLoop(config structs.Config, filesToWatch map[string]string) {
	for {
		updatedFiles, _ := symfony.GetWatchMap(config)
		if !reflect.DeepEqual(filesToWatch, updatedFiles) {
			start := time.Now()
			fmt.Println()
			fmt.Println(fmt.Sprintf(" > %s at %s > refreshing cache", color.New(color.FgHiYellow).Sprintf("Update detected"), color.New(color.FgGreen).Sprintf(start.Format("15:04:05"))))
			_, _ = symfony.CacheWarmup(config)
			end := time.Now()
			elapsed := end.Sub(start)
			fmt.Println(fmt.Sprintf(" > %s in %s", color.New(color.FgGreen).Sprintf("Done"), color.New(color.FgHiYellow).Sprintf("%s", FormatDuration(elapsed.Milliseconds()))))
			filesToWatch = updatedFiles
			fmt.Println(fmt.Sprintf(" > %s file(s) watched at %s", color.YellowString("%d", len(filesToWatch)), color.YellowString("%s", config.DirSymfonyProject)))
			fmt.Println(fmt.Sprintf(" > %s to stop watching or run %s %s.", color.GreenString("CTRL+C"), color.GreenString("kill -9"), color.GreenString("%d", os.Getpid())))
		} else {
			time.Sleep(config.SleepTime)
		}
	}
}

func FormatDuration(ms int64) string {
	const (
		msInSecond = 1000
		msInMinute = 60000
	)

	minutes := ms / msInMinute
	ms = ms % msInMinute
	seconds := ms / msInSecond
	ms = ms % msInSecond

	var result []string

	if minutes > 0 {
		result = append(result, fmt.Sprintf("%d minute(s)", minutes))
	}
	if seconds > 0 {
		result = append(result, fmt.Sprintf("%d second(s)", seconds))
	}
	if ms > 0 || len(result) == 0 {
		result = append(result, fmt.Sprintf("%d millisecond(s)", ms))
	}

	return strings.Join(result, ", ")
}

// main is the entry point of the program. It initializes the configuration, displays a Welcome message,
// parses command line arguments, checks for required parameters, sets configuration values based on the command line arguments,
// checks for the existence of Symfony console, retrieves the Symfony version, gets the files to watch,
// displays some information about the project, and enters the main loop to monitor and react to file changes.
func main() {
	var config structs.Config
	var err error
	config.Init()

	fmt.Println()

	env := flag.String("env", "dev", "pass --env=env to the symfony console (default: dev)")
	noDebug := flag.Bool("no-debug", false, "pass --no-debug to the symfony console (default: false)")
	clearCache := flag.Bool("cache", false, "clear cache instead of just warmup (default: false)")
	forceClearCache := flag.Bool("force", false, "force clear cache (rm -rf var/cache) (default: false)")
	exclude := flag.String("exclude", "", "comma-separated directories not to watch")
	vendors := flag.String("vendor", "", "comma-separated list of vendors to watch")

	pools := structs.NewCustomFlag()
	flag.Var(pools, "pools", "comma-separated list of pools to clear")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	clickableVersion := fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", GenerateVersionLink(version), version)
	fmt.Println(fmt.Sprintf(" > Version: %s", color.New(color.FgHiYellow).Sprintf(clickableVersion)))

	config.SymfonyEnv = *env
	config.ClearCache = *clearCache
	config.SymfonyDebug = !*noDebug

	if *forceClearCache {
		config.ClearCache = false
		config.ForceClearCache = true
	}

	if *vendors != "" {
		vendorList := ParseCommaSeparated(*vendors)
		if len(vendorList) > 0 {
			config.VendorWatch = true
			config.VendorList = vendorList
		}
	}

	if *exclude != "" {
		excludeDirs := ParseCommaSeparated(*exclude)
		config.DirsExclude = append(config.DirsExclude, excludeDirs...)
	}

	if pools.IsChanged() {
		config.PoolsProvided = true
		config.Pools = pools.Get()

		if len(config.Pools) == 0 {
			config.Pools = []string{"--all"}
		}
	}

	config.DirSymfonyProject, err = symfony.GetSymfonyProjectDir()

	if err != nil {
		PrintError(fmt.Errorf("project directory not found"))
		PrintError(err)
		os.Exit(1)
	}

	fmt.Println(" > Project directory: " + color.New(color.FgGreen).Sprintf(config.DirSymfonyProject))

	err = symfony.CheckSymfonyConsole(config)
	if err != nil {
		PrintError(fmt.Errorf("symfony console not found"))
		PrintError(err)
		os.Exit(1)
	}

	fmt.Println(" > Symfony console path: " + color.New(color.FgGreen).Sprintf(config.SymfonyConsolePath))

	out, err := symfony.Version(config)
	if err != nil {
		PrintError(fmt.Errorf("error while running the Symfony version command"))
		PrintError(err)
		os.Exit(1)
	}

	fmt.Println(" > Symfony env: " + color.New(color.FgGreen).Sprintf(strings.TrimSpace(fmt.Sprintf("%s", out))))

	start := time.Now()
	filesToWatch, _ := symfony.GetWatchMap(config)
	end := time.Now()
	elapsed := end.Sub(start)

	if len(filesToWatch) == 0 {
		PrintError(fmt.Errorf("no file to watch found"))
		os.Exit(0)
	}

	fmt.Println(fmt.Sprintf(" > %s file(s) watched at %s in %s", color.YellowString("%d", len(filesToWatch)), color.YellowString("%s", config.DirSymfonyProject), color.YellowString("%s", FormatDuration(elapsed.Milliseconds()))))
	fmt.Println(fmt.Sprintf(" > %s to stop watching or run %s %s.", color.GreenString("CTRL+C"), color.GreenString("kill -9"), color.GreenString("%d", os.Getpid())))

	MainLoop(config, filesToWatch)
}

// ParseCommaSeparated splits a comma-separated input string and returns an array of strings.
// If the input string is empty, it returns an empty array.
func ParseCommaSeparated(input string) []string {
	if input == "" {
		return []string{}
	}

	return strings.Split(input, ",")
}

// GenerateSeparator returns a string consisting of a specified number of em dashes.
// The length parameter determines the number of em dashes in the output string.
// The function uses the strings.Repeat function to repeat the em dash character.
func GenerateSeparator(length int) string {
	return strings.Repeat("—", length)
}

// IsSemanticVersion checks if a given version string is in the semantic version format.
func IsSemanticVersion(version string) bool {
	semVerPattern := `^v?\d+\.\d+\.\d+(-[a-zA-Z0-9]+(\.[a-zA-Z0-9]+)*)?$`
	matched, _ := regexp.MatchString(semVerPattern, version)

	return matched
}

// GenerateVersionLink generates a link for a given version string.
// If the version is a semantic version, it creates a link to the release tag on the repository.
// If the version is not a semantic version, it creates a link to the commit on the repository.
// It returns the generated link as a string.
func GenerateVersionLink(version string) string {
	if IsSemanticVersion(version) || "nightly" == version {
		return fmt.Sprintf("%s/releases/tag/%s", repository, version)
	}

	return fmt.Sprintf("%s/commit/%s", repository, version)
}

// PrintError prints an error message to the console. If the given error is nil, it does nothing.
// Otherwise, it formats the error message with a red warning symbol and prints it in the console.
func PrintError(err error) {
	if err == nil {
		return
	}

	fmt.Println(fmt.Sprintf("%s %s /!\\", color.New(color.FgHiRed).Sprint("/!\\"), err))
}
