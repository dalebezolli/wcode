package config

import (
	"os"
	"path/filepath"
	"strings"
)

const directory = "$HOME/.config/wcode"

func dir() string {
	return os.ExpandEnv(directory)
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
