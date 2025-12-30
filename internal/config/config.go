package config

import (
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	DownloadPath string
}

var defaultConfig *Config

func GetDefaultDownloadsPath() string {
	var downloadsPath string

	switch runtime.GOOS {
	case "windows":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "."
		}
		downloadsPath = filepath.Join(homeDir, "Downloads")
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "."
		}
		downloadsPath = filepath.Join(homeDir, "Downloads")
	case "linux":
		if xdgDir := os.Getenv("XDG_DOWNLOAD_DIR"); xdgDir != "" {
			downloadsPath = xdgDir
		} else {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "."
			}
			downloadsPath = filepath.Join(homeDir, "Downloads")
		}
	default:
		return "."
	}

	return downloadsPath
}

func GetConfig() *Config {
	if defaultConfig == nil {
		defaultConfig = &Config{
			DownloadPath: GetDefaultDownloadsPath(),
		}
	}
	return defaultConfig
}

func (c *Config) GetDownloadPath() string {
	return c.DownloadPath
}

func (c *Config) SetDownloadPath(path string) {
	c.DownloadPath = path
}
