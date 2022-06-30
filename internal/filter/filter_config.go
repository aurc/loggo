package filter

type Config struct {
	Name        string      `json:"name" yaml:"name"`
	FilterGroup filterGroup `json:"filterGroup" yaml:"filter-group"`
}
