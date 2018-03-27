package store

import (
	"strings"

	"github.com/mdevilliers/dtree"
)

type repo struct {
	nodes []dtree.Node
	edges []dtree.Edge
}

func InMemory(nodes []dtree.Node, edges []dtree.Edge) *repo {
	return &repo{nodes: nodes, edges: edges}
}

func (r *repo) FindNodes(name string) ([]dtree.Node, error) {

	matches := []dtree.Node{}

	for _, n := range r.nodes {

		if name == n.Name {
			matches = append(matches, n)
			continue
		}

		if strings.Contains(n.Name, name) {
			matches = append(matches, n)
			continue
		}
	}

	return matches, nil
}

// naughty
var seen = map[string]bool{}

//FromNode starts at this node and return all dependancies recursively
func (r *repo) FromNode(node dtree.Node) ([]dtree.Node, []dtree.Edge) {

	nodes := map[string]dtree.Node{}
	edges := []dtree.Edge{}

	for _, e := range r.edges {

		if e.Relationship == dtree.Dependancy && e.Source.Name == node.Name {

			nodes[e.Source.Name] = e.Source
			nodes[e.Target.Name] = e.Target

			edges = append(edges, e)

			_, f := seen[e.Target.Name]

			if !f {

				n2, e2 := r.FromNode(e.Target)

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
func (r *repo) ToNode(node dtree.Node) ([]dtree.Node, []dtree.Edge) {

	nodes := map[string]dtree.Node{}
	edges := []dtree.Edge{}

	for _, e := range r.edges {

		if e.Relationship == dtree.Dependant && e.Target.Name == node.Name {

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
