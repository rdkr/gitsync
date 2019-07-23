package cmd

import (
	"fmt"
	"strings"
	"sync"

	gitlab "github.com/xanzy/go-gitlab"
)

type syncer struct {
	groups   chan group
	projects chan project

	groupsWG, projectsWG, projectsSignalWG *sync.WaitGroup
	projectSignalOnce                      *sync.Once

	ui ui
}

func newSyncer(ui ui) syncer {
	return syncer{
		groups:            make(chan group),
		projects:          make(chan project),
		groupsWG:          new(sync.WaitGroup),
		projectsWG:        new(sync.WaitGroup),
		projectsSignalWG:  new(sync.WaitGroup),
		projectSignalOnce: new(sync.Once),
		ui: ui,
	}
}


type group interface {
	getGroups() []group
	getProjects() []project
	rootFullPath() string
	rootLocation() string
}

type project interface {
	getPath() string
	getURL() string
}

type gitlabGroup struct {
	client *gitlab.Client
	fullPath string
	location string
	*gitlab.Group
}

type gitlabProject struct {
	client *gitlab.Client
	group group
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
		result = append(result, gitlabGroup{g.client, g.rootFullPath(), g.rootLocation(), group})
	}

	return result
}

func (g gitlabGroup) getProjects() []project {
	var result []project

	projects , _, err := g.client.Groups.ListGroupProjects(g.ID, &gitlab.ListGroupProjectsOptions{
		Archived: gitlab.Bool(false),
	})
	if err != nil {
		panic(err)
	}

	for _, project := range projects {
		result = append(result, gitlabProject{g.client, g, project})
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

func (s syncer) recurseGroups() {
	for {

		parent, ok := <-s.groups
		if !ok {
			break
		}

		childGroups := parent.getGroups()

		for _, child := range childGroups {
			s.groupsWG.Add(1)
			go func(group group) {
				s.groups <- group
			}(child)
		}

		childProjects := parent.getProjects()

		for _, child := range childProjects {
			s.projectsWG.Add(1)
			s.projectSignalOnce.Do(func() { s.projectsSignalWG.Done() })
			go func(project project) {
				s.projects <- project
			}(child)
		}

		s.groupsWG.Done()
	}
}

func (s syncer) processProject() {
	for {

		project, ok := <-s.projects
		if !ok {
			break
		}

		s.ui.statusChan <- clone(project.getPath(), project.getURL())
		s.projectsWG.Done()
	}
}
