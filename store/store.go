package store

import (
	"github.com/mdevilliers/dtree"
)

type repo struct {
	nodes []dtree.Node
	edges []dtree.Edge
}

func InMemory(nodes []dtree.Node, edges []dtree.Edge) *repo {
	return &repo{nodes: nodes, edges: edges}
}

func (r *repo) All() ([]dtree.Node, []dtree.Edge) {
	return r.nodes, r.edges
}

// naughty
var seen = map[string]bool{}

//FromNode starts at this node and return all dependancies recursively
func (r *repo) FromNode(name string) ([]dtree.Node, []dtree.Edge) {

	nodes := map[string]dtree.Node{}
	edges := []dtree.Edge{}

	for _, e := range r.edges {

		if e.Relationship == dtree.Dependancy && e.Source.Name == name {

			nodes[e.Source.Name] = e.Source
			nodes[e.Target.Name] = e.Target

			edges = append(edges, e)

			_, f := seen[e.Target.Name]

			if !f {

				n2, e2 := r.FromNode(e.Target.Name)

				for _, n := range n2 {
					nodes[n.Name] = n
				}

				edges = append(edges, e2...)

				seen[e.Source.Name] = true
				seen[e.Target.Name] = true
			}
		}
	}
	return mapToArr(nodes), edges

}

// ToNode finds all dependancies looking at this node
func (r *repo) ToNode(name string) ([]dtree.Node, []dtree.Edge) {

	nodes := map[string]dtree.Node{}
	edges := []dtree.Edge{}

	for _, e := range r.edges {

		if e.Relationship == dtree.Dependant && e.Target.Name == name {

			nodes[e.Source.Name] = e.Source

			edges = append(edges, e)
		}
	}
	return mapToArr(nodes), edges
}

// GroupAll returns all dependancies grouped by version
func (r *repo) GroupAll() map[string][]dtree.Edge {

	ret := map[string][]dtree.Edge{}
	for _, e := range r.edges {

		if e.Relationship == dtree.Dependancy {
			n := e.Target.Name
			_, contains := ret[n]
			if !contains {
				ret[n] = []dtree.Edge{e}
			} else {
				ret[n] = append(ret[n], e)
			}
		}
	}
	return ret
}

// GroupOn returns dependancies for grouped by version for a module
func (r *repo) GroupOn(name string) map[string][]dtree.Edge {

	all := r.GroupAll()

	ret := map[string][]dtree.Edge{}

	for k, _ := range all {

		if k == name {
			ret[k] = all[k]
			break
		}

	}
	return ret
}

func mapToArr(m map[string]dtree.Node) []dtree.Node {

	arr := []dtree.Node{}

	for _, n := range m {
		arr = append(arr, n)
	}

	return arr
}
