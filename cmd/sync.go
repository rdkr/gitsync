package cmd

import (
	"sync"
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
		ui:                ui,
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
