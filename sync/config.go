package sync

type Config struct {
	Gitlab gitlabConfig `yaml:"gitlab"`
	Anon   anonConfig   `yaml:"anon"`
}

type gitlabConfig struct {
	Groups   []gitlabGroup `yaml:"groups"`
	Projects []Project     `yaml:"projects"`
	Token    string        `yaml:"token"`
}

type anonConfig struct {
	Projects []Project `yaml:"projects"`
}

type gitlabGroup struct {
	Group    int    `yaml:"group"`
	Location string `yaml:"location"`
}

type Project struct {
	URL      string `yaml:"url"`
	Location string `yaml:"location"`
	Token    string `yaml:"token"`
}
