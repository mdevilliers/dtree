package parsers

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/mdevilliers/dtree"
)

type dep struct{}

type depFile struct {
	Constraint []depConstraint
}

type depConstraint struct {
	Name    string
	Version string
}

func (d dep) Test(pth string) bool {
	pathToGoDepFile := path.Join(pth, "Gopkg.toml")
	return fileExists(pathToGoDepFile)
}

func (d dep) Parse(pth string) ([]dtree.Node, []dtree.Edge, error) {

	pathToGoDepFile := path.Join(pth, "Gopkg.toml")

	data, err := ioutil.ReadFile(pathToGoDepFile) // nolint: gosec

	if err != nil {
		return nil, nil, err
	}

	df := depFile{}
	_, err = toml.Decode(string(data), &df)

	if err != nil {
		return nil, nil, fmt.Errorf("cannot unmarshal data: %v", err)
	}

	parts := strings.Split(pth, "/")
	rootName := strings.Join(parts[len(parts)-3:], "/")

	root := dtree.RootNode(rootName)
	nodes := []dtree.Node{root}

	edges := []dtree.Edge{}

	for _, p := range df.Constraint {

		version := p.Version
		if version == "" {
			version = master
		}

		node := dtree.NewNode(p.Name, version)

		pair := dtree.NewDependency(root, node, version)

		nodes = append(nodes, node)
		edges = append(edges, pair...)
	}

	return nodes, edges, nil

}
