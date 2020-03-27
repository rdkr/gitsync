package concurrency

import (
	"sync"
)

type Group interface {
	GetGroups() []Group
	GetProjects() []Project
}

type User interface {
	GetProjects() []Project
}

type Project struct {
	URL      string `yaml:"url"`
	Location string `yaml:"location"`
	Token    string `yaml:"token"`
} // TODO move git specific stuff to sync module

type projectActionFunc func(Project) interface{}
type projectChanSenderFunc func(projectAction projectActionFunc, project Project)
type projectsChanCloserFunc func()

type manager struct {
	groups   chan Group
	projects chan Project

	groupsWG, groupsSignalWG, projectsWG, projectsSignalWG *sync.WaitGroup
	groupsSignalOnce, projectsSignalOnce                   *sync.Once

	projectAction projectActionFunc
}

func newManager(projectAction projectActionFunc) manager {
	return manager{
		groups:             make(chan Group),
		projects:           make(chan Project),
		groupsWG:           new(sync.WaitGroup),
		groupsSignalWG:     new(sync.WaitGroup),
		groupsSignalOnce:   new(sync.Once),
		projectsWG:         new(sync.WaitGroup),
		projectsSignalWG:   new(sync.WaitGroup),
		projectsSignalOnce: new(sync.Once),
		projectAction:      projectAction,
	}
}

func (m manager) start(groups []Group, projects []Project, projectsChanSender projectChanSenderFunc, projectsChanCloser projectsChanCloserFunc) {

	var wg sync.WaitGroup
	wg.Add(2)

	m.groupsWG.Add(1)
	m.groupsSignalWG.Add(1)

	m.projectsWG.Add(2)
	m.projectsSignalWG.Add(1)

	// groups manager goroutine
	go func() {

		// start some groups processors
		for w := 0; w < 25; w++ {
			go m.processGroups()
		}

		// wait for a signal indicating that we have a group to process
		m.groupsSignalWG.Wait()

		// wait to finish processing all groups before closing channel
		m.groupsWG.Wait()
		close(m.groups)

		// ensure we have processed all groups before stopping on projects
		m.projectsWG.Done()

		// stop the groups manager goroutine
		wg.Done()

	}()

	// projects manager goroutine
	go func() {

		// Start some projects processors
		for w := 0; w < 50; w++ {
			go m.processProject(projectsChanSender)
		}

		// wait for a signal that we have a Project to process
		m.projectsSignalWG.Wait()

		// wait to finish processing all projects before closing channel
		m.projectsWG.Wait()
		close(m.projects)

		// ensure we have processed all projects before stopping the UI
		projectsChanCloser()

		// stop the projects manager goroutine
		wg.Done()

	}()

	for _, g := range groups {
		m.groupsWG.Add(1)
		m.groups <- g
	}

	m.groupsWG.Done()
	m.groupsSignalOnce.Do(func() { m.groupsSignalWG.Done() })

	for _, p := range projects {
		m.projectsWG.Add(1)
		m.projectsSignalOnce.Do(func() { m.projectsSignalWG.Done() })

		go func(project Project) {
			m.projects <- project
		}(p)
	}

	m.projectsWG.Done()
	m.projectsSignalOnce.Do(func() { m.projectsSignalWG.Done() })

	wg.Wait()

}

func (m manager) processGroups() {
	for {

		parent, ok := <-m.groups
		if !ok {
			break
		}

		childGroups := parent.GetGroups()

		for _, child := range childGroups {
			m.groupsWG.Add(1)
			m.groupsSignalOnce.Do(func() { m.groupsSignalWG.Done() })
			go func(group Group) {
				m.groups <- group
			}(child)
		}

		childProjects := parent.GetProjects()

		for _, child := range childProjects {
			m.projectsWG.Add(1)
			m.projectsSignalOnce.Do(func() { m.projectsSignalWG.Done() })
			go func(project Project) {
				m.projects <- project
			}(child)
		}

		m.groupsWG.Done()
	}
}

func (m manager) processProject(projectsChanSender projectChanSenderFunc) {
	for {

		project, ok := <-m.projects
		if !ok {
			break
		}

		projectsChanSender(m.projectAction, project)
		m.projectsWG.Done()
	}
}

// ChannelMerger merges the output of n channel into one
func ChannelMerger(output chan<- interface{}, inputs ...<-chan interface{}) {

	var wg sync.WaitGroup

	for _, input := range inputs {
		wg.Add(1)

		go func(input <-chan interface{}) {
			for {
				value, ok := <-input
				if !ok {
					break
				}
				output <- value
			}
			wg.Done()
		}(input)
	}

	wg.Wait()
	close(output)
}
