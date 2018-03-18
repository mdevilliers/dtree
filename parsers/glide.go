package parsers

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/go-yaml/yaml"
	"github.com/mdevilliers/dtree"
)

var (
	relationshipStr = "relationship"
)

type glide struct{}

type glideFile struct {
	Package string        `yaml:"package"`
	Imports []glideImport `yaml:"import"`
}
type glideImport struct {
	Package string `json:"package"`
	Version string `json:"version"`
}

func Glide() glide {
	return glide{}
}

func (g glide) Parse(pth string) ([]dtree.Node, []dtree.Edge, error) {

	pathToGlideFile := path.Join(pth, "glide.yaml")

	data, err := ioutil.ReadFile(pathToGlideFile)

	if err != nil {
		return nil, nil, err
	}

	glideFile := glideFile{}

	err = yaml.Unmarshal([]byte(data), &glideFile)

	if err != nil {
		return nil, nil, fmt.Errorf("cannot unmarshal data: %v", err)
	}

	root := dtree.RootNode(glideFile.Package)
	nodes := []dtree.Node{root}

	edges := []dtree.Edge{}

	for _, p := range glideFile.Imports {

		version := p.Version
		if version == "" {
			version = "master"
		}

		node := dtree.NewNode(p.Package, version)

		pair := dtree.NewDependancy(root, node, version)

		nodes = append(nodes, node)
		edges = append(edges, pair...)
	}

	return nodes, edges, nil
}
