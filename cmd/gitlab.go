package cmd

import (
	"fmt"
	"strings"

	gitlab "github.com/xanzy/go-gitlab"
)

type gitlabGroupProvider struct {
	client   *gitlab.Client
	token    string
	fullPath string
	location string
	*gitlab.Group
}

func (g gitlabGroupProvider) getGroups() []group {
	var result []group

	groups, _, err := g.client.Groups.ListSubgroups(g.ID, &gitlab.ListSubgroupsOptions{
		AllAvailable: gitlab.Bool(true),
	})
	if err != nil {
		panic(err)
	}

	for _, group := range groups {
		result = append(result, gitlabGroupProvider{g.client, g.token, g.rootFullPath(), g.rootLocation(), group})
	}

	return result
}

func (g gitlabGroupProvider) getProjects() []project {
	var result []project

	projects, _, err := g.client.Groups.ListGroupProjects(g.ID, &gitlab.ListGroupProjectsOptions{
		Archived: gitlab.Bool(false),
	})
	if err != nil {
		panic(err)
	}

	for _, p := range projects {

		path := strings.ReplaceAll(p.PathWithNamespace, g.rootFullPath(), "")
		path = strings.TrimLeft(path, "/")
		path = fmt.Sprintf("%s/%s", g.rootLocation(), path)

		result = append(result, project{p.HTTPURLToRepo, path, g.token})
	}

	return result
}

func (g gitlabGroupProvider) rootFullPath() string {
	return g.fullPath
}

func (g gitlabGroupProvider) rootLocation() string {
	return g.location
}
