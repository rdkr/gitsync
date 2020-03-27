package concurrency_test

import (
	"sort"
	"testing"

	"github.com/go-test/deep"
	"github.com/rdkr/gitsync/concurrency"
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
	mockGetItemsFromCfg func() ([]concurrency.ProviderProcessor, []concurrency.Project)
}{
	{
		name: "NoGroupsNoProjects",
		mockGetItemsFromCfg: func() ([]concurrency.ProviderProcessor, []concurrency.Project) {
			var groups []concurrency.ProviderProcessor
			var projects []concurrency.Project
			return groups, projects
		},
	},
	{
		name: "NoGroupsAProject",
		mockGetItemsFromCfg: func() ([]concurrency.ProviderProcessor, []concurrency.Project) {
			var groups []concurrency.ProviderProcessor
			projects := []concurrency.Project{
				{URL: "a", Location: "b", Token: "c"},
			}
			return groups, projects
		},
	},
	{
		name: "EmptyGroupNoProject",
		mockGetItemsFromCfg: func() ([]concurrency.ProviderProcessor, []concurrency.Project) {
			groups := []concurrency.ProviderProcessor{
				&testGroupProvider{children: nil, projects: nil},
			}
			var projects []concurrency.Project
			return groups, projects
		},
	},
	{
		name: "EmptyGroupAProject",
		mockGetItemsFromCfg: func() ([]concurrency.ProviderProcessor, []concurrency.Project) {
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
		mockGetItemsFromCfg: func() ([]concurrency.ProviderProcessor, []concurrency.Project) {
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
		mockGetItemsFromCfg: func() ([]concurrency.ProviderProcessor, []concurrency.Project) {
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

			groups, projects := tc.mockGetItemsFromCfg()

			m := concurrency.NewGitlabManager(func(project concurrency.Project) interface{} {
				return nil
			})

			go m.Start(groups, projects)

			for {
				_, ok := <-m.ProjectChan
				if !ok {
					break
				}
			}
		})
	}
}

func TestChannelMerger(t *testing.T) {

	out := make(chan interface{})
	in1 := make(chan interface{})
	in2 := make(chan interface{})
	in3 := make(chan interface{})

	go func() {
		concurrency.ChannelMerger(out, in1, in2, in3)
	}()

	go func() {
		in1 <- 5
		in2 <- 6
		in1 <- 1
		in2 <- 4
		in3 <- 3
		in3 <- 2
		close(in1)
		close(in2)
		close(in3)
	}()

	var outCheck []int
	for i := 0; i < 6; i++ {
		value := <-out
		outCheck = append(outCheck, value.(int))
	}

	sort.Ints(outCheck)
	if diff := deep.Equal(outCheck, []int{1, 2, 3, 4, 5, 6}); diff != nil {
		t.Error(diff)
	}
}
