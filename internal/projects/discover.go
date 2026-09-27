package projects

import (
	"os"
	"path/filepath"
)

func DiscoverProjects(roots []string) ([]string, error) {
	directories := []string{}
	for _, root := range roots {
		entries, err := os.ReadDir(root)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			directories = append(directories, filepath.Join(root, entry.Name()))
		}
	}

	return directories, nil
}
