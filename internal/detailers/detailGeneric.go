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
