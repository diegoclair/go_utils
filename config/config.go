package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// envKeyReplacer replaces dots and dashes with underscores for environment variable mapping.
var envKeyReplacer = strings.NewReplacer(".", "_", "-", "_")

// Options configures how configuration is loaded.
type Options struct {
	// ConfigFilePath is the explicit path to the config file.
	// When set, SearchPaths, ConfigName, and ConfigType are ignored.
	ConfigFilePath string

	// ConfigName is the config file name without extension.
	// Default: "config".
	ConfigName string

	// ConfigType is the config file type (e.g., "toml", "yaml", "json").
	// Default: "toml".
	ConfigType string

	// SearchPaths are directories where the config file will be searched.
	// Used when ConfigFilePath is not set.
	SearchPaths []string

	// UseEnvName when true, appends the value of the ENV environment variable
	// to ConfigName with a dash separator.
	// For example, with ConfigName="config" and ENV="local", it resolves to "config-local".
	// If ENV is not set, defaults to "local".
	UseEnvName bool

	// WatchChanges enables automatic reloading when the config file changes.
	WatchChanges bool
}

// Load reads configuration from a file and environment variables into a new instance of T.
// T must be a struct type compatible with viper's mapstructure unmarshaling.
//
// The loading process:
//  1. Reads the config file from the specified path or search paths.
//  2. Overrides config values with matching environment variables.
//     Environment variable names are derived by uppercasing config keys
//     and replacing dots/dashes with underscores (e.g., "db.host-name" → "DB_HOST_NAME").
//  3. Unmarshals the final configuration into a new *T.
//  4. Optionally watches the config file for changes and reloads automatically.
func Load[T any](opts Options) (*T, error) {
	v := viper.New()
	v.AutomaticEnv()

	if opts.ConfigFilePath != "" {
		v.SetConfigFile(opts.ConfigFilePath)
	} else {
		name := opts.ConfigName
		if name == "" {
			name = "config"
		}

		if opts.UseEnvName {
			env := os.Getenv("ENV")
			if env == "" {
				env = "local"
			}
			name = fmt.Sprintf("%s-%s", name, env)
		}

		configType := opts.ConfigType
		if configType == "" {
			configType = "toml"
		}

		v.SetConfigName(name)
		v.SetConfigType(configType)

		for _, path := range opts.SearchPaths {
			v.AddConfigPath(path)
		}
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	overrideWithEnvVars(v)

	cfg := new(T)
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	if opts.WatchChanges {
		v.WatchConfig()
		v.OnConfigChange(func(in fsnotify.Event) {
			if in.Op == fsnotify.Write {
				if err := v.Unmarshal(cfg); err != nil {
					slog.Error("failed to unmarshal config file changes",
						slog.String("error", err.Error()),
					)
				}
			}
		})
	}

	return cfg, nil
}

// overrideWithEnvVars iterates over all config keys and overrides their values
// with matching environment variables. Keys are converted to uppercase with
// dots and dashes replaced by underscores.
func overrideWithEnvVars(v *viper.Viper) {
	for _, k := range v.AllKeys() {
		key := strings.ToUpper(envKeyReplacer.Replace(k))
		if envValue := os.Getenv(key); envValue != "" {
			v.Set(k, envValue)
		}
	}
}
