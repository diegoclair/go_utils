# config Package

## Description

This package provides a generic configuration loader built on top of [viper](https://github.com/spf13/viper). It reads configuration from files and environment variables, unmarshaling into any struct type using Go generics.

## Features

- **Generic loading** — `Load[T]` returns a typed `*T`, no type assertions needed.
- **Environment variable override** — config keys are automatically mapped to env vars (e.g., `db.host-name` → `DB_HOST_NAME`).
- **Multiple file formats** — supports TOML, YAML, JSON, and any format viper supports.
- **Environment-based file resolution** — optionally appends the `ENV` variable to the config name (e.g., `config-local.toml`, `config-production.toml`).
- **Hot-reload** — optionally watches the config file and reloads on changes.
- **No global state** — uses a new viper instance per call, avoiding conflicts.

## Installation

```go
import "github.com/diegoclair/go_utils/config"
```

## Usage

### Basic usage with explicit file path

```go
type Config struct {
    App AppConfig `mapstructure:"app"`
    DB  DBConfig  `mapstructure:"db"`
}

cfg, err := config.Load[Config](config.Options{
    ConfigFilePath: "/path/to/config.toml",
})
```

### Using search paths

```go
cfg, err := config.Load[Config](config.Options{
    ConfigName:  "config",
    ConfigType:  "toml",
    SearchPaths: []string{".", "../", "../../deployment"},
})
```

### Environment-based config file

With `ENV=production`, this resolves to `config-production.toml`:

```go
cfg, err := config.Load[Config](config.Options{
    UseEnvName:  true,
    SearchPaths: []string{"./deployment"},
})
```

If `ENV` is not set, it defaults to `local` (i.e., `config-local.toml`).

### With hot-reload

```go
cfg, err := config.Load[Config](config.Options{
    ConfigFilePath: "/path/to/config.toml",
    WatchChanges:   true,
})
// cfg is automatically updated when the file changes.
```

### Singleton pattern (project-level)

The package does not enforce singleton behavior. Wrap it with `sync.Once` in your project:

```go
var (
    cfg     *Config
    cfgErr  error
    once    sync.Once
)

func GetConfig() (*Config, error) {
    once.Do(func() {
        cfg, cfgErr = config.Load[Config](config.Options{
            ConfigFilePath: resolveConfigPath(),
            WatchChanges:   true,
        })
    })
    return cfg, cfgErr
}
```

## Environment Variable Override

All config keys are checked against environment variables. The mapping converts keys to uppercase and replaces dots (`.`) and dashes (`-`) with underscores (`_`):

| Config Key       | Environment Variable |
|------------------|---------------------|
| `app.port`       | `APP_PORT`          |
| `db.host`        | `DB_HOST`           |
| `db.db-name`     | `DB_DB_NAME`        |
| `app.auth.key`   | `APP_AUTH_KEY`      |

Environment variables take precedence over file values.
