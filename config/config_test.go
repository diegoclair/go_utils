package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testConfig struct {
	App appConfig `mapstructure:"app"`
	DB  dbConfig  `mapstructure:"db"`
}

type appConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
	Port        string `mapstructure:"port"`
}

type dbConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db-name"`
}

const testTomlContent = `
[app]
name = "test-app"
environment = "test"
port = "8080"

[db]
host = "localhost"
port = 5432
username = "admin"
password = "secret"
db-name = "testdb"
`

func writeTestConfigFile(t *testing.T, dir, filename, content string) string {
	t.Helper()

	path := filepath.Join(dir, filename)
	err := os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err)

	return path
}

func TestLoad(t *testing.T) {

	t.Run("should load config from explicit file path", func(t *testing.T) {
		dir := t.TempDir()
		filePath := writeTestConfigFile(t, dir, "config.toml", testTomlContent)

		cfg, err := Load[testConfig](Options{
			ConfigFilePath: filePath,
		})

		require.NoError(t, err)
		assert.Equal(t, "test-app", cfg.App.Name)
		assert.Equal(t, "test", cfg.App.Environment)
		assert.Equal(t, "8080", cfg.App.Port)
		assert.Equal(t, "localhost", cfg.DB.Host)
		assert.Equal(t, 5432, cfg.DB.Port)
		assert.Equal(t, "admin", cfg.DB.Username)
		assert.Equal(t, "secret", cfg.DB.Password)
		assert.Equal(t, "testdb", cfg.DB.DBName)
	})

	t.Run("should load config from search paths", func(t *testing.T) {
		dir := t.TempDir()
		writeTestConfigFile(t, dir, "myconfig.toml", testTomlContent)

		cfg, err := Load[testConfig](Options{
			ConfigName:  "myconfig",
			ConfigType:  "toml",
			SearchPaths: []string{dir},
		})

		require.NoError(t, err)
		assert.Equal(t, "test-app", cfg.App.Name)
		assert.Equal(t, "localhost", cfg.DB.Host)
	})

	t.Run("should use default config name and type when not specified", func(t *testing.T) {
		dir := t.TempDir()
		writeTestConfigFile(t, dir, "config.toml", testTomlContent)

		cfg, err := Load[testConfig](Options{
			SearchPaths: []string{dir},
		})

		require.NoError(t, err)
		assert.Equal(t, "test-app", cfg.App.Name)
	})

	t.Run("should append env name when UseEnvName is true", func(t *testing.T) {
		dir := t.TempDir()
		writeTestConfigFile(t, dir, "config-staging.toml", testTomlContent)

		t.Setenv("ENV", "staging")

		cfg, err := Load[testConfig](Options{
			UseEnvName:  true,
			SearchPaths: []string{dir},
		})

		require.NoError(t, err)
		assert.Equal(t, "test-app", cfg.App.Name)
	})

	t.Run("should default to local when UseEnvName is true and ENV is empty", func(t *testing.T) {
		dir := t.TempDir()
		writeTestConfigFile(t, dir, "config-local.toml", testTomlContent)

		t.Setenv("ENV", "")

		cfg, err := Load[testConfig](Options{
			UseEnvName:  true,
			SearchPaths: []string{dir},
		})

		require.NoError(t, err)
		assert.Equal(t, "test-app", cfg.App.Name)
	})

	t.Run("should override config values with environment variables", func(t *testing.T) {
		dir := t.TempDir()
		filePath := writeTestConfigFile(t, dir, "config.toml", testTomlContent)

		t.Setenv("APP_PORT", "9090")
		t.Setenv("DB_USERNAME", "env-user")

		cfg, err := Load[testConfig](Options{
			ConfigFilePath: filePath,
		})

		require.NoError(t, err)
		assert.Equal(t, "9090", cfg.App.Port)
		assert.Equal(t, "env-user", cfg.DB.Username)
		// Non-overridden values remain from file.
		assert.Equal(t, "test-app", cfg.App.Name)
	})

	t.Run("should override nested keys with dashes via env vars", func(t *testing.T) {
		dir := t.TempDir()
		filePath := writeTestConfigFile(t, dir, "config.toml", testTomlContent)

		t.Setenv("DB_DB_NAME", "overridden-db")

		cfg, err := Load[testConfig](Options{
			ConfigFilePath: filePath,
		})

		require.NoError(t, err)
		assert.Equal(t, "overridden-db", cfg.DB.DBName)
	})

	t.Run("should load yaml config", func(t *testing.T) {
		dir := t.TempDir()

		yamlContent := `
app:
  name: yaml-app
  environment: dev
  port: "3000"
db:
  host: db-host
  port: 3306
  username: root
  password: pass
  db-name: yamldb
`
		writeTestConfigFile(t, dir, "app.yaml", yamlContent)

		cfg, err := Load[testConfig](Options{
			ConfigName:  "app",
			ConfigType:  "yaml",
			SearchPaths: []string{dir},
		})

		require.NoError(t, err)
		assert.Equal(t, "yaml-app", cfg.App.Name)
		assert.Equal(t, "db-host", cfg.DB.Host)
		assert.Equal(t, 3306, cfg.DB.Port)
	})

	t.Run("should return error when config file is not found", func(t *testing.T) {
		cfg, err := Load[testConfig](Options{
			ConfigFilePath: "/nonexistent/path/config.toml",
		})

		assert.Nil(t, cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "reading config file")
	})

	t.Run("should return error when search paths have no matching file", func(t *testing.T) {
		dir := t.TempDir()

		cfg, err := Load[testConfig](Options{
			ConfigName:  "missing",
			SearchPaths: []string{dir},
		})

		assert.Nil(t, cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "reading config file")
	})
}
