package core

import (
	"os"
	"path/filepath"
)

const (
	VersionCacheFile = "versions.json"
)

func GetInstallPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".gdcli", "versions")
}