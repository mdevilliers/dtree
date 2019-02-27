package parsers

import (
	"errors"
	"os"

	"github.com/mdevilliers/dtree"
)

const (
	master = "master"
)

// Parser will examis the provided path and return a collection
// of Nodes and Edges or an error
type Parser interface {
	Parse(pth string) (dtree.Nodes, dtree.Edges, error)
}

// New returns a parser instance or an error
func New(path string) (Parser, error) {

	d := dep{}
	if d.Test(path) {
		return d, nil
	}

	g := glide{}

	if g.Test(path) {
		return g, nil
	}

	return nil, errors.New("parser not found")
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
