package cmd

type config struct {
	Gitlab gitlabConfig `yaml:"gitlab"`
	Anon   anonConfig   `yaml:"anon"`
}

type gitlabConfig struct {
	Groups   []gitlabGroup `yaml:"groups"`
	Projects []project     `yaml:"projects"`
}

type anonConfig struct {
	Projects []project `yaml:"projects"`
}

type gitlabGroup struct {
	Group    int    `yaml:"group"`
	Location string `yaml:"location"`
}

type project struct {
	URL      string `yaml:"url"`
	Location string `yaml:"location"`
	Token    string `yaml:"token"`
}
