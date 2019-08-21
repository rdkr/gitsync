package sync_test

import (
	"gitsync/sync"
	"testing"
)

func mockGitSync(p sync.Git, location string) sync.Status {
	return sync.Status{location, "", "", nil}
}

type testGroupProvider struct {
	children []sync.ProviderProcessor
	projects []sync.Project
}

func (g *testGroupProvider) GetGroups() []sync.ProviderProcessor {
	return g.children
}

func (g *testGroupProvider) GetProjects() []sync.Project {
	return g.projects
}

var concurrencyTests = []struct {
	name                string
	mockGetItemsFromCfg func(sync.Config) ([]sync.ProviderProcessor, []sync.Project)
}{
	{
		name: "NoGroupsNoProjects",
		mockGetItemsFromCfg: func(sync.Config) ([]sync.ProviderProcessor, []sync.Project) {
			var groups []sync.ProviderProcessor
			var projects []sync.Project
			return groups, projects
		},
	},
	{
		name: "NoGroupsAProject",
		mockGetItemsFromCfg: func(sync.Config) ([]sync.ProviderProcessor, []sync.Project) {
			var groups []sync.ProviderProcessor
			projects := []sync.Project{
				{URL: "a", Location: "b", Token: "c"},
			}
			return groups, projects
		},
	},
	{
		name: "EmptyGroupNoProject",
		mockGetItemsFromCfg: func(sync.Config) ([]sync.ProviderProcessor, []sync.Project) {
			groups := []sync.ProviderProcessor{
				&testGroupProvider{children: nil, projects: nil},
			}
			var projects []sync.Project
			return groups, projects
		},
	},
	{
		name: "EmptyGroupAProject",
		mockGetItemsFromCfg: func(sync.Config) ([]sync.ProviderProcessor, []sync.Project) {
			groups := []sync.ProviderProcessor{
				&testGroupProvider{children: nil, projects: nil},
			}
			projects := []sync.Project{
				{URL: "a", Location: "b", Token: "c"},
			}
			return groups, projects
		},
	},
	{
		name: "NestedGroupNoProject",
		mockGetItemsFromCfg: func(sync.Config) ([]sync.ProviderProcessor, []sync.Project) {
			groups := []sync.ProviderProcessor{
				&testGroupProvider{
					children: []sync.ProviderProcessor{
						&testGroupProvider{children: nil, projects: nil},
					},
					projects: []sync.Project{
						{URL: "a", Location: "b", Token: "c"},
					},
				},
			}
			var projects []sync.Project
			return groups, projects
		},
	},
	{
		name: "NestedGroupAProject",
		mockGetItemsFromCfg: func(sync.Config) ([]sync.ProviderProcessor, []sync.Project) {
			groups := []sync.ProviderProcessor{
				&testGroupProvider{
					children: []sync.ProviderProcessor{
						&testGroupProvider{children: nil, projects: nil},
					},
					projects: []sync.Project{
						{URL: "a", Location: "b", Token: "c"},
					},
				},
			}
			projects := []sync.Project{
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
			sync.NewConcurrencyManager(sync.Config{}, sync.NewUI(true, true), tc.mockGetItemsFromCfg, mockGitSync).Start()
		})
	}
}
