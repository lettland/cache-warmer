package structs

import "time"

// Symfony default parameters for Symfony/Flex.
const (
	ConsolePath     = "bin/console"
	ClearCache      = false
	Debug           = true
	Env             = "dev"
	VendorWatch     = false
	DirConfig       = "config"
	DirMigrations   = "migrations"
	DirSrc          = "src"
	DirTemplates    = "templates"
	DirTranslations = "translations"
	DirVendor       = "vendor"
	ForceClearCache = false
	PoolsProvided   = false
	SleepTime       = 30 * time.Millisecond // Watcher process sleep time
)

// DefaultExcludedDirs contains the directories that should be excluded by default.
var DefaultExcludedDirs = []string{".git", ".github", "node_modules"}

// Config holds all the parameters needed for the application. The YAML tags
// represent the keys in the Symfony custom config file, which will override
// these default values.
type Config struct {
	ClearCache             bool // Clear cache instead of only warmup
	DirMigrations          string
	DirSymfonyConfig       string        // Directory where configuration files are stored
	DirSymfonyProject      string        // The main Symfony project directory
	DirSymfonySrc          string        // Directory where source code is stored
	DirSymfonyTemplates    string        // Directory where template files are stored
	DirSymfonyTranslations string        // Directory where translation files are stored
	DirSymfonyVendor       string        // Directory where vendor code is stored
	DirsExclude            []string      // Directories to exclude from monitoring
	ForceClearCache        bool          // Force cache removal using rm -rf var/cache
	Pools                  []string      // List of pools to watch
	PoolsProvided          bool          // Whether the --pools flag was provided
	SleepTime              time.Duration // Sleep time between filesystem checks
	SymfonyConsolePath     string        // Relative path to the Symfony console
	SymfonyDebug           bool          // APP_DEBUG parameter
	SymfonyEnv             string        // APP_ENV parameter
	VendorList             []string      // List of specific vendor directories to watch
	VendorWatch            bool          // Whether to watch vendor directories
}

// Init initializes the Config object with default values.
func (obj *Config) Init() {
	obj.ClearCache = ClearCache
	obj.DirMigrations = DirMigrations
	obj.DirSymfonyConfig = DirConfig
	obj.DirSymfonySrc = DirSrc
	obj.DirSymfonyTemplates = DirTemplates
	obj.DirsExclude = DefaultExcludedDirs
	obj.ForceClearCache = ForceClearCache
	obj.Pools = []string{}
	obj.PoolsProvided = PoolsProvided
	obj.SleepTime = SleepTime
	obj.SymfonyConsolePath = ConsolePath
	obj.SymfonyDebug = Debug
	obj.SymfonyEnv = Env
	obj.DirSymfonyTranslations = DirTranslations
	obj.DirSymfonyVendor = DirVendor
	obj.VendorList = []string{}
	obj.VendorWatch = VendorWatch
}
