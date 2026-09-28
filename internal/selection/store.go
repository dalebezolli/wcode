package selection

import (
	"os"

	"github.com/dalebezolli/wcode/internal/config"
)

func Save(path string) error {
	storePath, err := config.GetSelectionStorePath()
	if err != nil {
		return err
	}
	return os.WriteFile(storePath, []byte(path), 0666)
}

func Load() (string, error) {
	storePath, err := config.GetSelectionStorePath()
	if err != nil {
		return "", err
	}
	selected, err := os.ReadFile(storePath)
	if err != nil {
		return "", err
	}
	return string(selected), nil
}
