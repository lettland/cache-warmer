package structs

import (
	"reflect"
	"testing"
)

func TestConfig_Init(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   Config
	}{
		{
			name:   "DefaultValuesCheck",
			config: Config{},
			want: Config{
				ClearCache:             ClearCache,
				DirMigrations:          DirMigrations,
				DirSymfonyConfig:       DirConfig,
				DirSymfonySrc:          DirSrc,
				DirSymfonyTemplates:    DirTemplates,
				DirsExclude:            DefaultExcludedDirs,
				ForceClearCache:        ForceClearCache,
				Pools:                  []string{},
				PoolsProvided:          PoolsProvided,
				SleepTime:              SleepTime,
				SymfonyConsolePath:     ConsolePath,
				SymfonyDebug:           Debug,
				SymfonyEnv:             Env,
				DirSymfonyTranslations: DirTranslations,
				DirSymfonyVendor:       DirVendor,
				VendorList:             []string{},
				VendorWatch:            VendorWatch,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.config.Init()
			if !reflect.DeepEqual(tt.config, tt.want) {
				t.Errorf("Config.Init() = %v, want %v", tt.config, tt.want)
			}
		})
	}
}
