package sync

import "github.com/rdkr/gitsync/concurrency"

type Config struct {
	Github githubConfig `yaml:"github"`
	Gitlab gitlabConfig `yaml:"gitlab"`
	Anon   anonConfig   `yaml:"anon"`
}

type githubConfig struct {
	// Groups   []gitlabGroup         `yaml:"groups"`
	// Projects []concurrency.Project `yaml:"projects"`
	Users []string `yaml:"users"`
	Token string   `yaml:"token"`
}

// type gitlabGroup struct {
// 	Group    int    `yaml:"group"`
// 	Location string `yaml:"location"`
// }

type gitlabConfig struct {
	Groups   []gitlabGroup         `yaml:"groups"`
	Projects []concurrency.Project `yaml:"projects"`
	Token    string                `yaml:"token"`
}

type gitlabGroup struct {
	Group    int    `yaml:"group"`
	Location string `yaml:"location"`
}

type anonConfig struct {
	Projects []concurrency.Project `yaml:"projects"`
}
