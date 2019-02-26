package dtree

import "strings"

// Node is a dependency
type Node struct {
	Name    string
	Labels  Attributes
	Version string
}

// Edge connects two Nodes
type Edge struct {
	Source       Node
	Target       Node
	Labels       Attributes
	Version      string
	Relationship relationship
}

// RootNode creates the root Node for the tree
func RootNode(name string) Node {
	return Node{
		Name:    name,
		Version: "unknown",
		Labels:  Attributes{},
	}

}

// NewNode initilises a Node with a name and a version
func NewNode(name, version string) Node {
	return Node{
		Name:    name,
		Version: version,
		Labels:  Attributes{},
	}
}

// NewDependency connects two nodes with some version metadata
func NewDependency(source, target Node, version string) []Edge {
	return []Edge{
		{
			Source:       source,
			Target:       target,
			Labels:       Attributes{},
			Version:      version,
			Relationship: Dependency,
		},
		{
			Source:       source,
			Target:       target,
			Version:      version,
			Labels:       Attributes{},
			Relationship: Dependant,
		},
	}
}

type relationship string

var (
	// Dependant is a 'required by' relationship
	Dependant = relationship("dependent")
	// Dependency is a 'requires' relationship
	Dependency = relationship("dependency")
)

// Attributes is a collection of objects indexed by a name
type Attributes map[string]interface{}

// Nodes is a collection of Node objects
type Nodes []Node

// Predicate returms Nodes whose Name property contains a value
// An empty string matches all Nodes
func (n Nodes) Predicate(value string) Nodes {

	if value == "" {
		return n
	}

	matching := Nodes{}

	for _, nn := range n {
		if strings.Contains(nn.Name, value) {
			matching = append(matching, nn)
		}
	}

	return matching
}

// Edges is a collection of Edge nodes
type Edges []Edge

// Predicate returms Edges whose Nodes has a Name property contains a value
// An empty string matches all Edges
func (e Edges) Predicate(value string) Edges {

	if value == "" {
		return e
	}

	matching := Edges{}

	for _, ee := range e {
		if strings.Contains(ee.Source.Name, value) && strings.Contains(ee.Target.Name, value) {
			matching = append(matching, ee)
		}
	}

	return matching
}
