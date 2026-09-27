package config

import (
	"os"
	"path/filepath"
	"strings"
)

const wcodeConfigDir = "wcode"

func dir() string {
	userConfigRootDir, _ := os.UserConfigDir()
	return os.ExpandEnv(filepath.Join(userConfigRootDir, wcodeConfigDir))
}

func EnsureDir() error {
	return os.MkdirAll(dir(), 0751)
}

func GetProjectRoots() []string {
	return strings.Split(os.Getenv("WCODE_PATHS"), ";")
}

func GetSelectionStorePath() string {
	return filepath.Join(dir(), "selection")
}
