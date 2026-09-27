package selection

import (
	"os"

	"github.com/dalebezolli/wcode/internal/config"
)

func Save(path string) error {
	return os.WriteFile(config.GetSelectionStorePath(), []byte(path), 0666)
}
