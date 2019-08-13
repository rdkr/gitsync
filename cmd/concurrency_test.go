package cmd_test

import (
	"gitsync/cmd"
	"testing"
)

func mockGitSync(p cmd.Git, location string) cmd.Status {
	return cmd.Status{location, "", "", nil}
}

type testGroupProvider struct {
	children []cmd.ProviderProcessor
	projects []cmd.Project
}

func (g *testGroupProvider) GetGroups() []cmd.ProviderProcessor {
	return g.children
}

func (g *testGroupProvider) GetProjects() []cmd.Project {
	return g.projects
}

var concurrencyTests = []struct {
	name                string
	mockGetItemsFromCfg func() ([]cmd.ProviderProcessor, []cmd.Project)
}{
	{
		name: "NoGroupsNoProjects",
		mockGetItemsFromCfg: func() ([]cmd.ProviderProcessor, []cmd.Project) {
			var groups []cmd.ProviderProcessor
			var projects []cmd.Project
			return groups, projects
		},
	},
	{
		name: "NoGroupsAProject",
		mockGetItemsFromCfg: func() ([]cmd.ProviderProcessor, []cmd.Project) {
			var groups []cmd.ProviderProcessor
			projects := []cmd.Project{
				{URL: "a", Location: "b", Token: "c"},
			}
			return groups, projects
		},
	},
	{
		name: "EmptyGroupNoProject",
		mockGetItemsFromCfg: func() ([]cmd.ProviderProcessor, []cmd.Project) {
			groups := []cmd.ProviderProcessor{
				&testGroupProvider{children: nil, projects: nil},
			}
			var projects []cmd.Project
			return groups, projects
		},
	},
	{
		name: "EmptyGroupAProject",
		mockGetItemsFromCfg: func() ([]cmd.ProviderProcessor, []cmd.Project) {
			groups := []cmd.ProviderProcessor{
				&testGroupProvider{children: nil, projects: nil},
			}
			projects := []cmd.Project{
				{URL: "a", Location: "b", Token: "c"},
			}
			return groups, projects
		},
	},
	{
		name: "NestedGroupNoProject",
		mockGetItemsFromCfg: func() ([]cmd.ProviderProcessor, []cmd.Project) {
			groups := []cmd.ProviderProcessor{
				&testGroupProvider{
					children: []cmd.ProviderProcessor{
						&testGroupProvider{children: nil, projects: nil},
					},
					projects: []cmd.Project{
						{URL: "a", Location: "b", Token: "c"},
					},
				},
			}
			var projects []cmd.Project
			return groups, projects
		},
	},
	{
		name: "NestedGroupAProject",
		mockGetItemsFromCfg: func() ([]cmd.ProviderProcessor, []cmd.Project) {
			groups := []cmd.ProviderProcessor{
				&testGroupProvider{
					children: []cmd.ProviderProcessor{
						&testGroupProvider{children: nil, projects: nil},
					},
					projects: []cmd.Project{
						{URL: "a", Location: "b", Token: "c"},
					},
				},
			}
			projects := []cmd.Project{
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
			cmd.NewConcurrencyManager(cmd.NewUI(true), tc.mockGetItemsFromCfg, mockGitSync).Start()
		})
	}
}
