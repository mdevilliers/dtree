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

// naughty
var seen = map[string]bool{}

func (r *repo) FromNode(name string) ([]dtree.Node, []dtree.Edge) {

	nodes := map[string]dtree.Node{}
	edges := []dtree.Edge{}

	// start at this node - return all dependancies and theirs ect
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

func mapToArr(m map[string]dtree.Node) []dtree.Node {

	arr := []dtree.Node{}

	for _, n := range m {
		arr = append(arr, n)
	}

	return arr
}
