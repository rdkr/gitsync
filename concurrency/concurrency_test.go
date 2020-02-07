package concurrency_test

import (
	"github.com/rdkr/gitsync/concurrency"
	"testing"
)

type testGroupProvider struct {
	children []concurrency.ProviderProcessor
	projects []concurrency.Project
}

func (g *testGroupProvider) GetGroups() []concurrency.ProviderProcessor {
	return g.children
}

func (g *testGroupProvider) GetProjects() []concurrency.Project {
	return g.projects
}

var concurrencyTests = []struct {
	name                string
	mockGetItemsFromCfg func(concurrency.Config) ([]concurrency.ProviderProcessor, []concurrency.Project)
}{
	{
		name: "NoGroupsNoProjects",
		mockGetItemsFromCfg: func(concurrency.Config) ([]concurrency.ProviderProcessor, []concurrency.Project) {
			var groups []concurrency.ProviderProcessor
			var projects []concurrency.Project
			return groups, projects
		},
	},
	{
		name: "NoGroupsAProject",
		mockGetItemsFromCfg: func(concurrency.Config) ([]concurrency.ProviderProcessor, []concurrency.Project) {
			var groups []concurrency.ProviderProcessor
			projects := []concurrency.Project{
				{URL: "a", Location: "b", Token: "c"},
			}
			return groups, projects
		},
	},
	{
		name: "EmptyGroupNoProject",
		mockGetItemsFromCfg: func(concurrency.Config) ([]concurrency.ProviderProcessor, []concurrency.Project) {
			groups := []concurrency.ProviderProcessor{
				&testGroupProvider{children: nil, projects: nil},
			}
			var projects []concurrency.Project
			return groups, projects
		},
	},
	{
		name: "EmptyGroupAProject",
		mockGetItemsFromCfg: func(concurrency.Config) ([]concurrency.ProviderProcessor, []concurrency.Project) {
			groups := []concurrency.ProviderProcessor{
				&testGroupProvider{children: nil, projects: nil},
			}
			projects := []concurrency.Project{
				{URL: "a", Location: "b", Token: "c"},
			}
			return groups, projects
		},
	},
	{
		name: "NestedGroupNoProject",
		mockGetItemsFromCfg: func(concurrency.Config) ([]concurrency.ProviderProcessor, []concurrency.Project) {
			groups := []concurrency.ProviderProcessor{
				&testGroupProvider{
					children: []concurrency.ProviderProcessor{
						&testGroupProvider{children: nil, projects: nil},
					},
					projects: []concurrency.Project{
						{URL: "a", Location: "b", Token: "c"},
					},
				},
			}
			var projects []concurrency.Project
			return groups, projects
		},
	},
	{
		name: "NestedGroupAProject",
		mockGetItemsFromCfg: func(concurrency.Config) ([]concurrency.ProviderProcessor, []concurrency.Project) {
			groups := []concurrency.ProviderProcessor{
				&testGroupProvider{
					children: []concurrency.ProviderProcessor{
						&testGroupProvider{children: nil, projects: nil},
					},
					projects: []concurrency.Project{
						{URL: "a", Location: "b", Token: "c"},
					},
				},
			}
			projects := []concurrency.Project{
				{URL: "a", Location: "b", Token: "c"},
			}
			return groups, projects
		},
	},
}

// TestConcurrency ensures combinations of cfg inputs can be processed and do not hang the programme
func TestConcurrency(t *testing.T) {
	for _, tc := range concurrencyTests {
		t.Run(tc.name, func(t *testing.T) {

			groups, projects := tc.mockGetItemsFromCfg(concurrency.Config{})

			m := concurrency.NewManager(concurrency.Config{}, func(project concurrency.Project) concurrency.Status {
				return concurrency.Status{
					Path:   "none",
					Status: 0,
					Output: "ok",
					Err:    nil,
				}
			})

			go m.Start(groups, projects)

			for {
				_, ok := <-m.StatusChan
				if !ok {
					break
				}
			}
		})
	}
}
