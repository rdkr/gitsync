package cmd

import (
	"fmt"
	"strings"

	gitlab "github.com/xanzy/go-gitlab"
)

type gitlabGroup struct {
	client   *gitlab.Client
	token    string
	fullPath string
	location string
	*gitlab.Group
}

type gitlabProject struct {
	client *gitlab.Client
	token  string
	group  group
	*gitlab.Project
}

func (g gitlabGroup) getGroups() []group {
	var result []group

	groups, _, err := g.client.Groups.ListSubgroups(g.ID, &gitlab.ListSubgroupsOptions{
		AllAvailable: gitlab.Bool(true),
	})
	if err != nil {
		panic(err)
	}

	for _, group := range groups {
		result = append(result, gitlabGroup{g.client, g.token, g.rootFullPath(), g.rootLocation(), group})
	}

	return result
}

func (g gitlabGroup) getProjects() []project {
	var result []project

	projects, _, err := g.client.Groups.ListGroupProjects(g.ID, &gitlab.ListGroupProjectsOptions{
		Archived: gitlab.Bool(false),
	})
	if err != nil {
		panic(err)
	}

	for _, project := range projects {
		result = append(result, gitlabProject{g.client, g.token, g, project})
	}

	return result
}

func (g gitlabGroup) rootFullPath() string {
	return g.fullPath
}

func (g gitlabGroup) rootLocation() string {
	return g.location
}

func (p gitlabProject) getPath() string {
	path := strings.ReplaceAll(p.PathWithNamespace, p.group.rootFullPath(), "")
	path = strings.TrimLeft(path, "/")
	path = fmt.Sprintf("%s/%s", p.group.rootLocation(), path)
	return path
}

func (p gitlabProject) getURL() string {
	return p.HTTPURLToRepo
}

func (p gitlabProject) getToken() string {
	return p.token
}
