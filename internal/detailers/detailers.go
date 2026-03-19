package detailers

type Details struct {
	Title string
	Path  string

	Rest map[string]string
}

type Detailer interface {
	GetDetails(path string) Details
	GetRestOrder() []string
}

type DetailerType string

const (
	DetailerTypeClassic DetailerType = "classic"
	DetailerTypeGit     DetailerType = "git"
)

func NewDetailer(detailerType DetailerType) Detailer {
	switch detailerType {
	case DetailerTypeClassic:
		return DetailerClassic{}

	case DetailerTypeGit:
		if !IsDetailerGitAvailable() {
			break
		}
		return DetailerGit{}
	}

	return DetailerClassic{}
}
