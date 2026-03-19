package detailers

import (
	"strings"
)

type DetailerClassic struct {
}

func (d DetailerClassic) GetDetails(path string) Details {
	var title string
	splitPath := strings.Split(path, "/")
	title = splitPath[len(splitPath)-1]

	return Details{
		Title: title,
		Path:  path,
	}
}

func (d DetailerClassic) GetRestOrder() []string {
	return []string{}
}
