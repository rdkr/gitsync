package sync

import "github.com/rdkr/gitsync/concurrency"

type Config struct {
	Gitlab gitlabConfig `yaml:"gitlab"`
	Anon   anonConfig   `yaml:"anon"`
}

type gitlabConfig struct {
	Groups   []gitlabGroup         `yaml:"groups"`
	Projects []concurrency.Project `yaml:"projects"`
	Token    string                `yaml:"token"`
}

type anonConfig struct {
	Projects []concurrency.Project `yaml:"projects"`
}

type gitlabGroup struct {
	Group    int    `yaml:"group"`
	Location string `yaml:"location"`
}
