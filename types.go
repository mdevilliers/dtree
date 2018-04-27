package dtree

import "strings"

type Node struct {
	Name    string
	Labels  Attributes
	Version string
}

type Edge struct {
	Source       Node
	Target       Node
	Labels       Attributes
	Version      string
	Relationship relationship
}

func RootNode(name string) Node {
	return Node{
		Name:    name,
		Version: "unknown",
		Labels:  Attributes{},
	}

}

func NewNode(name, version string) Node {
	return Node{
		Name:    name,
		Version: version,
		Labels:  Attributes{},
	}
}

func NewDependancy(source, target Node, version string) []Edge {
	return []Edge{
		{
			Source:       source,
			Target:       target,
			Labels:       Attributes{},
			Version:      version,
			Relationship: Dependancy,
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
	Dependant  = relationship("dependent")
	Dependancy = relationship("dependancy")
)

type Attributes map[string]interface{}

type Nodes []Node

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

type Edges []Edge

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
