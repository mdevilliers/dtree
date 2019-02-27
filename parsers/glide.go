package parsers

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/go-yaml/yaml"
	"github.com/mdevilliers/dtree"
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

func (glide) Test(pth string) bool {
	pathToGlideFile := path.Join(pth, "glide.yaml")
	return fileExists(pathToGlideFile)
}

func (g glide) Parse(pth string) (dtree.Nodes, dtree.Edges, error) {

	pathToGlideFile := path.Join(pth, "glide.yaml")

	data, err := ioutil.ReadFile(pathToGlideFile) // nolint: gosec

	if err != nil {
		return nil, nil, err
	}

	gf := glideFile{}

	err = yaml.Unmarshal(data, &gf)

	if err != nil {
		return nil, nil, fmt.Errorf("cannot unmarshal data: %v", err)
	}

	root := dtree.RootNode(gf.Package)
	nodes := []dtree.Node{root}

	edges := []dtree.Edge{}

	for _, p := range gf.Imports {

		version := p.Version
		if version == "" {
			version = master
		}

		node := dtree.NewNode(p.Package, version)

		pair := dtree.NewDependency(root, node, version)

		nodes = append(nodes, node)
		edges = append(edges, pair...)
	}

	return nodes, edges, nil
}
