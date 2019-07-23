package cmd

type gitlabGroupConfig struct {
	Group    int    `yaml:"group"`
	Location string `yaml:"location"`
}

type gitlabConfig struct {
	Groups   []gitlabGroupConfig `yaml:"groups"`
	Projects []projectConfig     `yaml:"projects"`
}

type projectConfig struct {
	URL      string `yaml:"url"`
	Location string `yaml:"location"`
	Token    string `yaml:"token"`
}

type anonConfig struct {
	Projects []projectConfig `yaml:"projects"`
}

type config struct {
	Gitlab gitlabConfig `yaml:"gitlab"`
	Anon   anonConfig   `yaml:"anon"`
}
