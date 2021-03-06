package store

import (
	"strings"

	"github.com/mdevilliers/dtree"
)

type repo struct {
	nodes []dtree.Node
	edges []dtree.Edge
	seen  map[string]bool
}

// InMemory initilises an ephemeral store of source code information
func InMemory(nodes dtree.Nodes, edges dtree.Edges) *repo { // nolint: golint
	return &repo{nodes: nodes, edges: edges, seen: map[string]bool{}}
}

func (r *repo) FindNodes(name string) (dtree.Nodes, error) {

	matches := dtree.Nodes{}

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

//FromNode starts at this node and return all dependencies recursively
func (r *repo) FromNode(node dtree.Node) (dtree.Nodes, dtree.Edges) {

	nodes := map[string]dtree.Node{}
	edges := []dtree.Edge{}

	for _, e := range r.edges {

		if e.Relationship == dtree.Dependency && e.Source.Name == node.Name {

			nodes[e.Source.Name] = e.Source
			nodes[e.Target.Name] = e.Target

			edges = append(edges, e)

			_, f := r.seen[e.Target.Name]

			if !f {

				n2, e2 := r.FromNode(e.Target)

				for _, n := range n2 {
					nodes[n.Name] = n
				}

				edges = append(edges, e2...)

				r.seen[e.Source.Name] = true
				r.seen[e.Target.Name] = true
			}
		}
	}
	return mapToArr(nodes), edges

}

// ToNode finds all dependencies looking at this node
func (r *repo) ToNode(node dtree.Node) (dtree.Nodes, dtree.Edges) {

	nodes := map[string]dtree.Node{}
	edges := dtree.Edges{}

	for _, e := range r.edges {

		if e.Relationship == dtree.Dependant && e.Target.Name == node.Name {

			nodes[e.Source.Name] = e.Source

			edges = append(edges, e)
		}
	}
	return mapToArr(nodes), edges
}

func mapToArr(m map[string]dtree.Node) dtree.Nodes {

	arr := dtree.Nodes{}

	for _, n := range m {
		arr = append(arr, n)
	}

	return arr
}
