package config

import (
	"os"
	"path/filepath"
	"strings"
)

const wcodeConfigDir = "wcode"

func GetConfigRoot() (string, error) {
	userConfigRootDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return os.ExpandEnv(filepath.Join(userConfigRootDir, wcodeConfigDir)), nil
}

func EnsureDir() error {
	root, err := GetConfigRoot()
	if err != nil {
		return err
	}
	return os.MkdirAll(root, 0751)
}

func GetProjectRoots() []string {
	return strings.Split(os.Getenv("WCODE_PATHS"), ";")
}

func GetSelectionStorePath() (string, error) {
	root, err := GetConfigRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "selection"), nil
}
