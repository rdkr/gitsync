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

type Org interface {
	GetProjectsByOrg() []Project
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
	users    chan User
	orgs     chan Org
	groups   chan Group
	projects chan Project

	usersWG, usersSignalWG, orgsWG, orgsSignalWG, groupsWG, groupsSignalWG, projectsWG, projectsSignalWG *sync.WaitGroup
	usersSignalOnce, orgsSignalOnce, groupsSignalOnce, projectsSignalOnce                                *sync.Once

	projectAction projectActionFunc
}

func newManager(projectAction projectActionFunc) manager {
	return manager{
		users:              make(chan User),
		orgs:               make(chan Org),
		groups:             make(chan Group),
		projects:           make(chan Project),
		usersWG:            new(sync.WaitGroup),
		usersSignalWG:      new(sync.WaitGroup),
		usersSignalOnce:    new(sync.Once),
		orgsWG:             new(sync.WaitGroup),
		orgsSignalWG:       new(sync.WaitGroup),
		orgsSignalOnce:     new(sync.Once),
		groupsWG:           new(sync.WaitGroup),
		groupsSignalWG:     new(sync.WaitGroup),
		groupsSignalOnce:   new(sync.Once),
		projectsWG:         new(sync.WaitGroup),
		projectsSignalWG:   new(sync.WaitGroup),
		projectsSignalOnce: new(sync.Once),
		projectAction:      projectAction,
	}
}

func (m manager) start(users []User, orgs []Org, groups []Group, projects []Project, projectsChanSender projectChanSenderFunc, projectsChanCloser projectsChanCloserFunc) {

	var wg sync.WaitGroup
	wg.Add(4)

	m.usersWG.Add(1)
	m.usersSignalWG.Add(1)

	m.orgsWG.Add(1)
	m.orgsSignalWG.Add(1)

	m.groupsWG.Add(1)
	m.groupsSignalWG.Add(1)

	m.projectsWG.Add(4)
	m.projectsSignalWG.Add(1)

	// groups manager goroutine
	go func() {

		// start some groups processors
		for w := 0; w < 5; w++ {
			go m.processUsers()
		}

		// wait for a signal indicating that we have a group to process
		m.usersSignalWG.Wait()

		// wait to finish processing all groups before closing channel
		m.usersWG.Wait()
		close(m.users)

		// ensure we have processed all groups before stopping on projects
		m.projectsWG.Done()

		// stop the groups manager goroutine
		wg.Done()

	}()

	//org go routine
	go func() {

		// start some groups processors
		for w := 0; w < 5; w++ {
			go m.processOrgs()
		}

		// wait for a signal indicating that we have a group to process
		m.orgsSignalWG.Wait()

		// wait to finish processing all groups before closing channel
		m.orgsWG.Wait()
		close(m.orgs)

		// ensure we have processed all groups before stopping on projects
		m.projectsWG.Done()

		// stop the groups manager goroutine
		wg.Done()

	}()

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

	for _, u := range users {
		m.usersWG.Add(1)
		m.users <- u
	}
	m.usersWG.Done()
	m.usersSignalOnce.Do(func() { m.usersSignalWG.Done() })

	for _, o := range orgs {
		m.orgsWG.Add(1)
		m.orgs <- o
	}
	m.orgsWG.Done()
	m.orgsSignalOnce.Do(func() { m.orgsSignalWG.Done() })

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

func (m manager) processOrgs() {
	for {

		org, ok := <-m.orgs
		if !ok {
			break
		}

		childProjects := org.GetProjectsByOrg()
		for _, child := range childProjects {
			m.projectsWG.Add(1)
			m.projectsSignalOnce.Do(func() { m.projectsSignalWG.Done() })
			go func(project Project) {
				m.projects <- project
			}(child)
		}

		m.orgsWG.Done()
	}
}

func (m manager) processUsers() {
	for {

		user, ok := <-m.users
		if !ok {
			break
		}

		childProjects := user.GetProjects()
		for _, child := range childProjects {
			m.projectsWG.Add(1)
			m.projectsSignalOnce.Do(func() { m.projectsSignalWG.Done() })
			go func(project Project) {
				m.projects <- project
			}(child)
		}

		m.usersWG.Done()
	}
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
