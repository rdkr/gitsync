package concurrency

import (
	"sync"
)

const (
	StatusError = iota
	StatusCloned
	StatusFetched
	StatusUpToDate
)

type Status struct {
	Path   string
	Status int
	Output string
	Err    error
}

type ProviderProcessor interface {
	GetGroups() []ProviderProcessor
	GetProjects() []Project
}

type Manager struct {
	cfg Config

	groups   chan ProviderProcessor
	projects chan Project

	groupsWG, groupsSignalWG, projectsWG, projectsSignalWG *sync.WaitGroup
	groupsSignalOnce, projectsSignalOnce                   *sync.Once

	StatusChan    chan Status
	projectAction func(Project) Status
}

func NewManager(cfg Config, projectAction func(Project) Status) Manager {
	return Manager{
		cfg:                cfg,
		groups:             make(chan ProviderProcessor),
		projects:           make(chan Project),
		groupsWG:           new(sync.WaitGroup),
		groupsSignalWG:     new(sync.WaitGroup),
		groupsSignalOnce:   new(sync.Once),
		projectsWG:         new(sync.WaitGroup),
		projectsSignalWG:   new(sync.WaitGroup),
		projectsSignalOnce: new(sync.Once),
		StatusChan:         make(chan Status),
		projectAction:      projectAction,
	}
}

func (cm Manager) Start(groups []ProviderProcessor, projects []Project) {

	var wg sync.WaitGroup
	wg.Add(2)

	cm.groupsWG.Add(1)
	cm.groupsSignalWG.Add(1)

	cm.projectsWG.Add(2)
	cm.projectsSignalWG.Add(1)

	// groups manager goroutine
	go func() {

		// start some groups processors
		for w := 0; w < 10; w++ {
			go cm.processGroups()
		}

		// wait for a signal indicating that we have a group to process
		cm.groupsSignalWG.Wait()

		// wait to finish processing all groups before closing channel
		cm.groupsWG.Wait()
		close(cm.groups)

		// ensure we have processed all groups before stopping on projects
		cm.projectsWG.Done()

		// stop the groups manager goroutine
		wg.Done()

	}()

	// projects manager goroutine
	go func() {

		// Start some projects processors
		for w := 0; w < 20; w++ {
			go cm.processProject()
		}

		// wait for a signal that we have a Project to process
		cm.projectsSignalWG.Wait()

		// wait to finish processing all projects before closing channel
		cm.projectsWG.Wait()
		close(cm.projects)

		// ensure we have processed all projects before stopping the UI
		close(cm.StatusChan)

		// stop the projects manager goroutine
		wg.Done()

	}()

	for _, g := range groups {
		cm.groupsWG.Add(1)
		cm.groups <- g
	}

	cm.groupsWG.Done()
	cm.groupsSignalOnce.Do(func() { cm.groupsSignalWG.Done() })

	for _, p := range projects {
		cm.projectsWG.Add(1)
		cm.projectsSignalOnce.Do(func() { cm.projectsSignalWG.Done() })

		go func(project Project) {
			cm.projects <- project
		}(p)
	}

	cm.projectsWG.Done()
	cm.projectsSignalOnce.Do(func() { cm.projectsSignalWG.Done() })

	wg.Wait()

}

func (cm Manager) processGroups() {
	for {

		parent, ok := <-cm.groups
		if !ok {
			break
		}

		childGroups := parent.GetGroups()

		for _, child := range childGroups {
			cm.groupsWG.Add(1)
			cm.groupsSignalOnce.Do(func() { cm.groupsSignalWG.Done() })
			go func(group ProviderProcessor) {
				cm.groups <- group
			}(child)
		}

		childProjects := parent.GetProjects()

		for _, child := range childProjects {
			cm.projectsWG.Add(1)
			cm.projectsSignalOnce.Do(func() { cm.projectsSignalWG.Done() })
			go func(project Project) {
				cm.projects <- project
			}(child)
		}

		cm.groupsWG.Done()
	}
}

func (cm Manager) processProject() {
	for {

		project, ok := <-cm.projects
		if !ok {
			break
		}

		cm.StatusChan <- cm.projectAction(project)
		cm.projectsWG.Done()
	}
}
